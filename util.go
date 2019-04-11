package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/google/logger"
)

func removeFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func renameDB(src, dest string) (err error) {
	_, err = os.Stat(src)
	if nil == err {
		if err = os.Rename(src, dest); err != nil {
			return err
		}
		return nil
	} else if os.IsNotExist(err) {
		return nil
	}
	return err
}

func renameBitmarkdDB() (finalErr error) {
	mainnetBlockDB := nodeDataDirMainnet + "/" + blockLevelDB
	// XXX: Does not know how to handle error yet; records its error now
	if err := renameDB(mainnetBlockDB, mainnetBlockDB+oldDBPostfix); err != nil {
		if !os.IsNotExist(err) {
			finalErr = ErrCombind(finalErr, err)
		}
	}

	mainnetIndexDB := nodeDataDirMainnet + "/" + indexLevelDB
	if err := renameDB(mainnetIndexDB, mainnetIndexDB+oldDBPostfix); err != nil {
		if !os.IsNotExist(err) {
			finalErr = ErrCombind(finalErr, err)
		}
	}

	testnetBlockDB := nodeDataDirTestnet + "/" + blockLevelDB
	if err := renameDB(testnetBlockDB, testnetBlockDB+oldDBPostfix); err != nil {
		if !os.IsNotExist(err) {
			finalErr = ErrCombind(finalErr, err)
		}
	}

	testnetIndexDB := nodeDataDirTestnet + "/" + indexLevelDB
	if err := renameDB(testnetIndexDB, testnetIndexDB+oldDBPostfix); err != nil {
		if !os.IsNotExist(err) {
			finalErr = ErrCombind(finalErr, err)
		}
	}
	return finalErr
}

func builDefaultVolumSrcBaseDir() (string, error) {
	homeDir := os.Getenv("USER_NODE_BASE_DIR")
	if 0 == len(homeDir) {
		return "", ErrorUserNodeDirEnv
	}
	return homeDir, nil
}

func recoverBitmarkdDB() (finalErr error) {
	mainnetBlockDB := nodeDataDirMainnet + "/" + blockLevelDB
	if err := renameDB(mainnetBlockDB+oldDBPostfix, mainnetBlockDB); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}

	mainnetIndexDB := nodeDataDirMainnet + "/" + indexLevelDB
	if err := renameDB(mainnetIndexDB+oldDBPostfix, mainnetIndexDB); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}

	testnetBlockDB := nodeDataDirTestnet + "/" + blockLevelDB
	if err := renameDB(testnetBlockDB+oldDBPostfix, testnetBlockDB); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}

	testnetIndexDB := nodeDataDirTestnet + "/" + indexLevelDB
	if err := renameDB(testnetIndexDB+oldDBPostfix, testnetIndexDB); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}
	return finalErr

}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
