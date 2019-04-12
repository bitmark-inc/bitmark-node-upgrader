package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tuserHomeDir() string {
	return "/home/pieceofr2"
}

var mockUserCorrectDir = map[string]string{
	"baseDir":        filepath.Join(tuserHomeDir(), "bitmark-node-data-test"),
	"nodeDBDir":      filepath.Join(tuserHomeDir(), "bitmark-node-data-test", "db"),
	"mainnetDataDir": filepath.Join(tuserHomeDir(), "bitmark-node-data-test", "data"),
	"testnetDataDir": filepath.Join(tuserHomeDir(), "bitmark-node-data-test", "data-test"),
	"mainnetLogDir":  filepath.Join(tuserHomeDir(), "bitmark-node-data-test", "log"),
	"testnetLogDir":  filepath.Join(tuserHomeDir(), "bitmark-node-data-test", "log-test"),
}

var mockDockerCorrectDir = map[string]string{
	"baseDir":            "/.config/bitmark-node",
	"nodeDBDir":          filepath.Join("/.config/bitmark-node", "db"),
	"mainnetDataDir":     filepath.Join("/.config/bitmark-node", "bitmarkd/bitmark", "data"),
	"testnetDataDir":     filepath.Join("/.config/bitmark-node", "bitmarkd/testing", "data"),
	"mainnetLogDir":      filepath.Join("/.config/bitmark-node", "bitmarkd/bitmark", "log"),
	"testnetLogDir":      filepath.Join("/.config/bitmark-node", "bitmarkd/testing", "log"),
	"mainnetBlockDBPath": filepath.Join("/.config/bitmark-node", "bitmarkd/bitmark", "data", "bitmark-blocks.leveldb"),
	"mainnetIndexDBPath": filepath.Join("/.config/bitmark-node", "bitmarkd/bitmark", "data", "bitmark-index.leveldb"),
	"mainnetZipFilePath": filepath.Join("/.config/bitmark-node", "bitmarkd/bitmark", "data", "snapshot.zip"),
	"testnetBlockDBPath": filepath.Join("/.config/bitmark-node", "bitmarkd/testing", "data", "bitmark-blocks.leveldb"),
	"testnetIndexDBPath": filepath.Join("/.config/bitmark-node", "bitmarkd/testing", "data", "bitmark-index.leveldb"),
	"testnetZipFilePath": filepath.Join("/.config/bitmark-node", "bitmarkd/testing", "data", "snapshot.zip"),
}

func TestUserPath(t *testing.T) {
	baseDir, err := builDefaultVolumSrcBaseDir()
	assert.NoError(t, err, fmt.Sprintf("Unable to get base directory"))

	userPath := UserPath{
		//	BaseDir:        "/home/pieceofr2/bitmark-node-data-test",
		BaseDir:        baseDir,
		NodeDBDir:      "db",
		MainnetDataDir: "data",
		TestnetDataDir: "data-test",
		MainnetLogDir:  "log",
		TestnetLogDir:  "log-test",
	}

	assert.Equal(t, mockUserCorrectDir["baseDir"], userPath.BaseDir, "baseDir not match")
	assert.Equal(t, mockUserCorrectDir["nodeDBDir"], userPath.GetNodeDBPath(), "nodeDBDir not match")
	assert.Equal(t, mockUserCorrectDir["mainnetDataDir"], userPath.GetDataPath(userPath.GetMainnet()), "mainnet Data dir not match")
	assert.Equal(t, mockUserCorrectDir["testnetDataDir"], userPath.GetDataPath(userPath.GetTestnet()), "testnet Data dir not match")
	assert.Equal(t, mockUserCorrectDir["mainnetLogDir"], userPath.GetLogPath(userPath.GetMainnet()), "mainnet Data dir not match")
	assert.Equal(t, mockUserCorrectDir["testnetLogDir"], userPath.GetLogPath(userPath.GetTestnet()), "testnet Data dir not match")

}

func TestDockerPath(t *testing.T) {

	dockerPath := DockerPath{
		BaseDir:             "/.config/bitmark-node",
		NodeDBDir:           "db",
		MainnetDataDir:      "bitmarkd/bitmark/data",
		TestnetDataDir:      "bitmarkd/testing/data",
		MainnetLogDir:       "bitmarkd/bitmark/log",
		TestnetLogDir:       "bitmarkd/testing/log",
		OldContainerPostfix: ".old",
		OldDatabasePostfix:  ".old",
		BlockDBDirName:      "bitmark-blocks.leveldb",
		IndexDBDirName:      "bitmark-index.leveldb",
		UpdateDBZipName:     "snapshot.zip",
	}

	assert.Equal(t, mockDockerCorrectDir["baseDir"], dockerPath.BaseDir, "baseDir not match")
	assert.Equal(t, mockDockerCorrectDir["nodeDBDir"], dockerPath.GetNodeDBPath(), "nodeDBDir not match")
	assert.Equal(t, mockDockerCorrectDir["mainnetDataDir"], dockerPath.GetDataPath(dockerPath.GetMainnet()), "mainnet Data dir not match")
	assert.Equal(t, mockDockerCorrectDir["testnetDataDir"], dockerPath.GetDataPath(dockerPath.GetTestnet()), "testnet Data dir not match")
	assert.Equal(t, mockDockerCorrectDir["mainnetLogDir"], dockerPath.GetLogPath(dockerPath.GetMainnet()), "mainnet Data dir not match")
	assert.Equal(t, mockDockerCorrectDir["testnetLogDir"], dockerPath.GetLogPath(dockerPath.GetTestnet()), "testnet Data dir not match")

	assert.Equal(t, mockDockerCorrectDir["mainnetBlockDBPath"], dockerPath.GetBlockDBPath(dockerPath.GetMainnet()), "mainnetBlockDBPath not match")
	assert.Equal(t, mockDockerCorrectDir["mainnetIndexDBPath"], dockerPath.GetIndexDBPath(dockerPath.GetMainnet()), "mainnetIndexDBPath not match")
	assert.Equal(t, mockDockerCorrectDir["mainnetZipFilePath"], dockerPath.GetUpdateDBZipFilePath(dockerPath.GetMainnet()), "mainnetZipFilePath not match")

	assert.Equal(t, mockDockerCorrectDir["testnetBlockDBPath"], dockerPath.GetBlockDBPath(dockerPath.GetTestnet()), "testnetBlockDBPath not match")
	assert.Equal(t, mockDockerCorrectDir["testnetIndexDBPath"], dockerPath.GetIndexDBPath(dockerPath.GetTestnet()), "testnetIndexDBPath not match")
	assert.Equal(t, mockDockerCorrectDir["testnetZipFilePath"], dockerPath.GetUpdateDBZipFilePath(dockerPath.GetTestnet()), "testnetZipFilePath not match")
}
