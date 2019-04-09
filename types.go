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
type CreateConfig struct {
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

// ChaindataUpdater to get information
type ChaindataUpdater interface {
	IsUpdated() bool
	GetCurrentDBVersion() (int, error)
	GetLatestChainInfo() (*LatestChain, error)
}

// ChaindataUpdaterConfig Config interface
type ChaindataUpdaterConfig interface {
	GetConfig() ChaindataUpdaterConfig
}

// DBUpdaterHTTPS use Http to get Information
type DBUpdaterHTTPS struct {
	// Endpoint and Path
	LatestChainInfoEndpoint string
	CurrentDBPath           string
	CurrentDataPath         string
	CurrentDataTestPath     string
	ZipSourcePath           string
	ZipDestinationPath      string
	// Status
	CurrentDBVer int
	Latest       LatestChain
}

// DBUpdaterHTTPSConfig for configuraion of UpdateDB
type DBUpdaterHTTPSConfig struct {
	APIEndpoint         string
	CurrentDBPath       string
	CurrentDataPath     string
	CurrentDataTestPath string
	ZipSourcePath       string
	ZipDestinationPath  string
}

// GetConfig Get DBUpdaterHTTPSConfig itself
func (d DBUpdaterHTTPSConfig) GetConfig() ChaindataUpdaterConfig {
	return d
}

// LatestChain latest database info
type LatestChain struct {
	Created     string `json:"created"`
	Version     string `json:"version"`
	BlockHeight int    `json:"blockheight"`
	DataURL     string `json:"dataurl"`
}

// GetVerion return the Version
func (i *LatestChain) GetVerion() (int, error) {
	n, err := strconv.ParseInt(i.Version, 0, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
