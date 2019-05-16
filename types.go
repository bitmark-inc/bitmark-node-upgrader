package main

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
)

// NodeAction bitmarkd API type
type NodeAction int

const (
	bitmarkdStart NodeAction = iota
	bitmarkdStop
	recorderdStart
	recorderdStop
)

// BitmarkdServResp API response for bitmarkd service
type BitmarkdServResp struct {
	Message string `json:"msg"`
	OK      int    `json:"ok"`
}

// RecorderdServResp API response for bitmarkd service
type RecorderdServResp struct {
	Message string `json:"msg"`
	OK      int    `json:"ok"`
}

// UserPath is uesed to store user file path
type UserPath struct {
	BaseDir        string
	NodeDBDir      string
	MainnetDataDir string
	TestnetDataDir string
	MainnetLogDir  string
	TestnetLogDir  string
}

// DockerPath is uesed to store docker file path
type DockerPath struct {
	BaseDir             string
	NodeDBDir           string
	MainnetDataDir      string
	TestnetDataDir      string
	MainnetLogDir       string
	TestnetLogDir       string
	OldContainerPostfix string
	OldDatabasePostfix  string
	BlockDBDirName      string
	IndexDBDirName      string
	UpdateDBZipName     string
}

// NodeWatcher main data structure of service
type NodeWatcher struct {
	DockerClient     *dockerClient.Client
	BackgroundContex context.Context
	Repo             string
	ImageName        string
	ContainerName    string
	Postfix          string
}

// CreateContainerConfig collect configs to create a container
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
	IsStartFromGenesis() bool
	IsForceUpdate() bool
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
	ForceUpdate     bool   `json:"forceupdate"`
	Version         string `json:"version"`
	FromGenesis     bool   `json:"fromgenesis"`
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
