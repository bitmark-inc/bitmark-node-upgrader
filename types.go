package main

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
)

// NodeWatcher main data structure of service
type NodeWatcher struct {
	DockerClient     *dockerClient.Client
	BackgroundContex context.Context
	Repo             string
	ImageName        string
	ContainerName    string
	Postfix          string
}

// CreateConfig collect configs to create a container
type CreateContainerConfig struct {
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

// DBUpdater to get information
type DBUpdater interface {
	IsUpdated() (main bool, test bool)
	GetCurrentDBVersion() (mainnet int, testbet int, err error)
	GetLatestChain() (*LatestChain, error)
}

// DBUpdaterConfig Config interface
type DBUpdaterConfig interface {
	GetConfig() DBUpdaterConfig
}

// DBUpdaterHTTPS use Http to get Information
type DBUpdaterHTTPS struct {
	// Endpoint and Path
	LatestChainInfoEndpoint string
	// Mainnet
	CurrentDBPath      string
	ZipSourcePath      string
	ZipDestinationPath string
	// Mainnet Status
	CurrentDBVer int
	// Testnet
	CurrentDBTestPath      string
	CurrentDataTestPath    string
	ZipSourceTestPath      string
	ZipDestinationTestPath string
	// Testnet Status
	CurrentTestDBVer int
	Latest           LatestChain
}

// DBUpdaterHTTPSConfig for configuraion of UpdateDB
type DBUpdaterHTTPSConfig struct {
	APIEndpoint string
	// Mainnet
	CurrentDBPath      string
	ZipSourcePath      string
	ZipDestinationPath string
	// Testnet
	CurrentDBTestPath      string
	CurrentDataTestPath    string
	ZipSourceTestPath      string
	ZipDestinationTestPath string
}

// GetConfig Get DBUpdaterHTTPSConfig itself
func (d DBUpdaterHTTPSConfig) GetConfig() DBUpdaterConfig {
	return d
}

// LatestChain latest database info
type LatestChain struct {
	Created         string `json:"created"`
	Version         string `json:"version"`
	BlockHeight     int    `json:"blockheight"`
	DataURL         string `json:"dataurl"`
	TestVersion     string `json:"testversion"`
	TestBlockHeight int64  `json:"testblockheight"`
	TestDataURL     string `json:"testdataurl"`
}

// GetVerion return the Version
func (i *LatestChain) GetVerion() (int, error) {
	n, err := strconv.ParseInt(i.Version, 0, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// GetVerionTestnet return the Version
func (i *LatestChain) GetVerionTestnet() (int, error) {
	n, err := strconv.ParseInt(i.TestVersion, 0, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
