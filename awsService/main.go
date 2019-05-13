package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

const (
	bucketRegion = "ap-northeast-1"
	bucket       = "bitmark-node-update"
	versionFile  = "latestDBVers.json"
)

// LatestChain latest database info
type LatestChain struct {
	Created         string `json:"created"`
	ForceUpdate     string `json:"forceupdate"`
	Version         string `json:"version"`
	BlockHeight     int64  `json:"blockheight"`
	FromGenesis     bool   `json:"fromgenesis"`
	DataURL         string `json:"dataurl"`
	TestVersion     string `json:"testversion"`
	TestBlockHeight int64  `json:"testblockheight"`
	TestDataURL     string `json:"testdataurl"`
}

func main() {
	lambda.Start(LambdaHandler)
}

// LambdaHandler to handle lambda request
func LambdaHandler() (*LatestChain, error) {
	ret, err := getLatestDBInfo()
	if err != nil && ret == nil {
		return nil, errors.New("get Last Database Info error")
	}

	return ret, nil
}

func getLatestDBInfo() (*LatestChain, error) {
	fmt.Println("testForUpdater start ... ")
	var latestchain LatestChain
	srv, err := NewS3Storage(bucket, bucketRegion)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	data, err := srv.GetData(versionFile)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	err = json.Unmarshal(data, &latestchain)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return nil, err
	}
	fmt.Println("data:", string(data))

	return &latestchain, nil
}

func (i *LatestChain) getVerion() (int64, error) {
	n, err := strconv.ParseInt(i.Version, 0, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
