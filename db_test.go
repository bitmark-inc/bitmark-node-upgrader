package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Copy folder of bitmark-node-data or bitmark-node-data-test into bitmark-node-data-unit-test for testing
// and change permission to avoid permission issue when not in docker
var LevelDBPath = filepath.Join(userHomeDir(), "bitmark-node-data-test")
var blockLevelDBPath = filepath.Join(userHomeDir(), "bitmark-node-data-test", "data", "bitmark-index.leveldb")

const (
	minVer      = 0x200
	dataEndpoit = "https://0u08da25ba.execute-api.ap-northeast-1.amazonaws.com/v1/chaindata"
)

func TestReadLevelDB(t *testing.T) {
	dbUpdater := DBUpdaterHTTPS{LatestChainInfoEndpoint: dataEndpoit, CurrentDBPath: blockLevelDBPath}
	ver, err := dbUpdater.GetCurrentDBVersion()
	assert.NoError(t, err, fmt.Sprintf("TestReadLevelDB:%s", err))
	assert.Equal(t, minVer, ver, fmt.Errorf("version is not equal").Error())
}

func TestGetLatestChainInfo(t *testing.T) {
	dbUpdater := DBUpdaterHTTPS{LatestChainInfoEndpoint: dataEndpoit}
	chainInfo, err := dbUpdater.GetLatestChainInfo()
	assert.NoError(t, err, fmt.Sprintf("TestGetLatestChainInfo:%s", err))
	assert.NotNil(t, chainInfo, fmt.Errorf("TestGetLatestChainInfo: chainInfo is nil"))
}

func TestSetDBUpdaterReady(t *testing.T) {
	config := DBUpdaterHTTPSConfig{APIEndpoint: dataEndpoit, CurrentDBPath: blockLevelDBPath}
	dbUpdater, err := SetDBUpdaterReady(config)
	assert.NoError(t, err, fmt.Sprintf("TestSetDBUpdaterReady:%s", err))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).LatestChainInfoEndpoint), 0, fmt.Sprintf("TestSetDBUpdaterReady: No LatestChainInfoEndpoint"))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).CurrentDBPath), 0, fmt.Sprintf("TestSetDBUpdaterReady: No CurrentDBPath"))
	assert.NotEqual(t, dbUpdater.(*DBUpdaterHTTPS).CurrentDBVer, 0, fmt.Sprintf("TestSetDBUpdaterReady: No CurrentDBPath CurrentDBVer:%d", dbUpdater.(*DBUpdaterHTTPS).CurrentDBVer))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).Latest.Version), 0, fmt.Sprintf("TestSetDBUpdaterReady: No Latest Chain Version"))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).Latest.DataURL), 0, fmt.Sprintf("TestSetDBUpdaterReady: No Latest DataURL"))
}

/*
func TestSetDownloadfile(t *testing.T) {
	zippath := filepath.Join(LevelDBPath, "test.zip")
	config := DBUpdaterHTTPSConfig{APIEndpoint: dataEndpoit,
		CurrentDBPath:   blockLevelDBPath,
		CurrentDataPath: LevelDBPath,
		ZipPath:         zippath}

	dbUpdater, err := SetDBUpdaterReady(config)
	fmt.Println(dbUpdater)
	assert.NoError(t, err, fmt.Sprintf("TestSetDownloadfile:%s", err))
	err = dbUpdater.(*DBUpdaterHTTPS).downloadfile()
	assert.NoError(t, err, fmt.Sprintf("TestSetDownloadfile:%s", err))
}
*/
func TestUpdateToLatestDB(t *testing.T) {
	baseDir, err := builDefaultVolumSrcBaseDir()
	assert.NoError(t, err, fmt.Sprintf("TestUpdateToLatestDB:baseDir:%s", err))

	zipsource := filepath.Join(baseDir, "data", "snapshot.zip")
	zipdestination := filepath.Join(baseDir, "data")

	config := DBUpdaterHTTPSConfig{APIEndpoint: dataEndpoit,
		CurrentDBPath:      blockLevelDBPath,
		CurrentDataPath:    LevelDBPath,
		ZipSourcePath:      zipsource,
		ZipDestinationPath: zipdestination,
	}
	dbUpdater, err := SetDBUpdaterReady(config)
	assert.NoError(t, err, fmt.Sprintf("TestUpdateToLatestDB:SetDBUpdaterReady:%s", err))
	err = dbUpdater.(*DBUpdaterHTTPS).UpdateToLatestDB()
	assert.NoError(t, err, fmt.Sprintf("TestUpdateToLatestDB:UpdateToLatestDB:%s", err))

}

/*
	baseDir, err := builDefaultVolumSrcBaseDir()
	if err != nil {
		return err
	}
	//dest := filepath.Join("/home/pieceofr2", "bitmark-node-data-unit-test")
	dest := filepath.Join(baseDir, "data")
*/
