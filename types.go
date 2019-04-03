package main

import (
	"context"
	"os"
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

// RemoteInfoServer to get information
type RemoteInfoServer interface {
	GetVersion() int
}

// RemoteDBServer is to get chain database
type RemoteDBServer interface {
	Download() os.File
}

// HTTPRemote use Http to get Information
type RemoteHTTPS3 struct {
	InfoEndpoint string
	APILatestVer string
	DBEndpoint   string
}

// DBInfo latest database info
type DBInfo struct {
	Created     string `json:"created"`
	Version     string `json:"version"`
	BlockHeight int64  `json:"blockheight"`
	DataURL     string `json:"dataurl"`
}

func (i *DBInfo) getVerion() (int64, error) {
	n, err := strconv.ParseInt(i.Version, 0, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
