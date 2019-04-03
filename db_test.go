package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var blockLevelDBPath = userHomeDir() + "/bitmark-node-data-test" + "/data/" + "bitmark-index.leveldb"

const (
	minVer      = 0x200
	dataEndpoit = "https://0u08da25ba.execute-api.ap-northeast-1.amazonaws.com/v1/chaindata"
)

func TestReadLevelDB(t *testing.T) {
	ver, err := readBmrkdDBVers(blockLevelDBPath)
	assert.NoError(t, err, fmt.Errorf("TestReadLevelDB: %s", err))
	fmt.Printf("leveldb verion:%d minVer:%d\n", ver, minVer)
	assert.Equal(t, minVer, ver, fmt.Errorf("version is not equal").Error())
}
