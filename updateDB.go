package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/google/logger"
	"github.com/syndtr/goleveldb/leveldb"
	ldb_opt "github.com/syndtr/goleveldb/leveldb/opt"
)

var versionKey = []byte{0x00, 'V', 'E', 'R', 'S', 'I', 'O', 'N'}

// SetDBUpdaterReady is to setup the specific type of RemoteLatestChainFetcher and RemoteDBDownloader
func SetDBUpdaterReady(conf ChaindataUpdaterConfig) (ChaindataUpdater, error) {
	updaterConfig := conf.(DBUpdaterHTTPSConfig).GetConfig()
	httpS3Updater := &DBUpdaterHTTPS{
		LatestChainInfoEndpoint: updaterConfig.(DBUpdaterHTTPSConfig).APIEndpoint,
		CurrentDBPath:           updaterConfig.(DBUpdaterHTTPSConfig).CurrentDBPath,
		CurrentDataPath:         updaterConfig.(DBUpdaterHTTPSConfig).CurrentDataPath,
		CurrentDataTestPath:     updaterConfig.(DBUpdaterHTTPSConfig).CurrentDataTestPath,
		ZipSourcePath:           updaterConfig.(DBUpdaterHTTPSConfig).ZipSourcePath,
		ZipDestinationPath:      updaterConfig.(DBUpdaterHTTPSConfig).ZipDestinationPath,
	}
	// get the currentDBVersion
	_, err := httpS3Updater.GetCurrentDBVersion()
	if err != nil {
		return httpS3Updater, err
	}
	latest, err := httpS3Updater.GetLatestChainInfo()
	if err != nil {
		return httpS3Updater, err
	}
	if latest != nil {
		httpS3Updater.Latest = *latest
	}

	return httpS3Updater, nil
}

// GetCurrentDBVersion get current chainData version
func (r *DBUpdaterHTTPS) GetCurrentDBVersion() (int, error) {

	opt := &ldb_opt.Options{
		ErrorIfExist:   false,
		ErrorIfMissing: true,
		ReadOnly:       true,
	}

	db, err := leveldb.OpenFile(r.CurrentDBPath, opt)
	if nil != err {
		return 0, err
	}

	versionValue, err := db.Get(versionKey, nil)
	if leveldb.ErrNotFound == err {
		return 0, nil
	} else if nil != err {
		return 0, err
	}
	if 4 != len(versionValue) {
		db.Close()
		log.Errorf("incompatible database version length: expected: %d  actual: %d", 4, len(versionValue))
		return 0, ErrorIncompatibleVersionLength
	}
	version := int(binary.BigEndian.Uint32(versionValue))
	r.CurrentDBVer = version
	db.Close()
	return version, nil
}

// IsUpdated is current databse updated
func (r *DBUpdaterHTTPS) IsUpdated() bool {
	if r.CurrentDBVer != 0 {
		latestVer, err := r.Latest.GetVerion()
		if err != nil {
			return false
		}
		if latestVer != r.CurrentDBVer {
			return false
		}
	}
	return true
}

// GetLatestChainInfo to get latestChainInfo from Retmote
func (r *DBUpdaterHTTPS) GetLatestChainInfo() (*LatestChain, error) {
	resp, err := http.Get(r.LatestChainInfoEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var latestChain LatestChain
	err = json.Unmarshal(body, &latestChain)
	if err != nil {
		return nil, err
	}
	return &latestChain, err
}

// UpdateToLatestDB Download latest and update the local database
func (r *DBUpdaterHTTPS) UpdateToLatestDB() error {
	/*
		if r.IsUpdated() {
			log.Info("UpdateToLatestDB IsUpdated")
			return nil
		}
	*/
	err := r.downloadfile()
	if err != nil {
		return err
	}
	err = renameBitmarkdDB()
	if err != nil {
		return err
	}
	err = unzip(r.ZipSourcePath, r.ZipDestinationPath)
	if err != nil {
		recoverErr := recoverBitmarkdDB()
		return ErrCombind(err, recoverErr)
	}
	fmt.Println("UpdateToLatestDB Successful")
	err = removeFile(r.ZipSourcePath) // nice to have so does not return error even it has error
	log.Warning("UpdateToLatestDB:remove zip file error:", err)
	return nil
}

func (r *DBUpdaterHTTPS) downloadfile() error {
	resp, err := http.Get(r.Latest.DataURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	if 0 == len(r.ZipSourcePath) {
		baseDir, err := builDefaultVolumSrcBaseDir()
		if err != nil {
			return err
		}
		r.ZipSourcePath = filepath.Join(baseDir, "levelDB.zip")
	}
	zipfile, err := os.Create(r.ZipSourcePath)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Write the body to file
	_, err = io.Copy(zipfile, resp.Body)
	if err != nil {
		return nil
	}

	return err
}
