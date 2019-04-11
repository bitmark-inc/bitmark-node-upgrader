package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Copy folder of bitmark-node-data or bitmark-node-data-test into bitmark-node-data-unit-test for testing
// and change permission to avoid permission issue when not in docker
var appBaseDir = filepath.Join(userHomeDir(), "bitmark-node-data-test")
var blockLevelDBPath = filepath.Join(userHomeDir(), "bitmark-node-data-test", "data", "bitmark-index.leveldb")

const (
	minVer      = 0x200
	dataEndpoit = "https://0u08da25ba.execute-api.ap-northeast-1.amazonaws.com/v1/chaindata"
)

func TestReadLevelDB(t *testing.T) {
	dbUpdater := DBUpdaterHTTPS{LatestChainInfoEndpoint: dataEndpoit, CurrentDBPath: blockLevelDBPath}
	ver, _, err := dbUpdater.GetCurrentDBVersion()
	// TODO: for testnet
	if !os.IsNotExist(err) {
		assert.NoError(t, err, fmt.Sprintf("TestReadLevelDB:%s", err))
	}
	if ver != 0 {
		assert.Equal(t, minVer, ver, fmt.Errorf("version is not equal").Error())
	}
}

func TestGetLatestChain(t *testing.T) {
	dbUpdater := DBUpdaterHTTPS{LatestChainInfoEndpoint: dataEndpoit}
	chainInfo, err := dbUpdater.GetLatestChain()
	assert.NoError(t, err, fmt.Sprintf("TestGetLatestChainInfo:%s", err))
	assert.NotNil(t, chainInfo, fmt.Errorf("TestGetLatestChainInfo: chainInfo is nil"))
}

func TestSetDBUpdaterReady(t *testing.T) {
	config := DBUpdaterHTTPSConfig{APIEndpoint: dataEndpoit, CurrentDBPath: blockLevelDBPath}
	dbUpdater, err := SetDBUpdaterReady(config)
	assert.NoError(t, err, fmt.Sprintf("TestSetDBUpdaterReady:%s", err))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).LatestChainInfoEndpoint), 0, fmt.Sprintf("TestSetDBUpdaterReady: No LatestChainInfoEndpoint"))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).CurrentDBPath), 0, fmt.Sprintf("TestSetDBUpdaterReady: No CurrentDBPath"))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).Latest.Version), 0, fmt.Sprintf("TestSetDBUpdaterReady: No Latest Chain Version"))
	assert.NotEqual(t, len(dbUpdater.(*DBUpdaterHTTPS).Latest.DataURL), 0, fmt.Sprintf("TestSetDBUpdaterReady: No Latest DataURL"))
}

func TestUpdateToLatestDB(t *testing.T) {
	zipsource := filepath.Join(appBaseDir, "data", "snapshot.zip")
	zipdestination := filepath.Join(appBaseDir, "data")

	config := DBUpdaterHTTPSConfig{
		APIEndpoint:        dataEndpoit,
		CurrentDBPath:      blockLevelDBPath,
		ZipSourcePath:      zipsource,
		ZipDestinationPath: zipdestination,
	}
	fmt.Println(config)
	dbUpdater, err := SetDBUpdaterReady(config)
	assert.NoError(t, err, fmt.Sprintf("TestUpdateToLatestDB:SetDBUpdaterReady:%s", err))
	err = dbUpdater.(*DBUpdaterHTTPS).UpdateToLatestDB()
	assert.NoError(t, err, fmt.Sprintf("TestUpdateToLatestDB:UpdateToLatestDB:%s", err))
	// TODO: for testnet

}
