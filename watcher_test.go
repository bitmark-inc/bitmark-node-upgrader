package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/client"
	log "github.com/google/logger"
	"github.com/stretchr/testify/assert"
)

var mockData *MockData

type MockData struct {
	Watcher    NodeWatcher
	userPath   UserPath
	dockerPath DockerPath
	Chain      string
	Env        map[string]string
}

func TestMain(m *testing.M) {
	logfile, err := os.OpenFile("unittest.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.Init("bitmark-node-updater-log", true, false, logfile)
	mockData = &MockData{}
	mockData.init()
	if err := cleanContainers(); err != nil {
		fmt.Println("Setup Error:", err)
		panic("cleanContainer Error")
	}
	os.Exit(m.Run())
}

func TestPullImage(t *testing.T) {
	watcher := mockData.getWatcher()
	_, err := watcher.pullImage()
	assert.NoError(t, err, ErrorImagePull.Error())
}

func TestStartContainer(t *testing.T) {
	watcher := mockData.getWatcher()
	newContainerConfig, err := getDefaultConfig(watcher)
	assert.NoError(t, err, ErrorConfigCreateNew.Error())
	newContainer, err := watcher.DockerClient.ContainerCreate(watcher.BackgroundContex, newContainerConfig.Config,
		newContainerConfig.HostConfig, nil, watcher.ContainerName)
	assert.NoError(t, err, ErrorContainerCreate.Error())

	err = watcher.startContainer(newContainer.ID)
	assert.NoError(t, err, ErrorContainerStart.Error())
}

func TestStopContainer(t *testing.T) {
	watcher := mockData.getWatcher()
	containers, _ := watcher.getContainersWithImage()
	err := watcher.stopContainers(containers, 10*time.Second)
	assert.NoError(t, err, ErrorContainerStop.Error())
}

func TestRenameContainer(t *testing.T) {
	watcher := mockData.getWatcher()
	containers, _ := watcher.getContainersWithImage()
	container := watcher.getNamedContainer(containers)
	assert.NotNil(t, container, 0, "No container to stop")
	watcher.renameContainer(container)
	container, err := watcher.getOldContainer()
	assert.NotNil(t, container, "Rename Container fail")
	assert.NoError(t, err, "get old container fail")

}
func (mock *MockData) init() error {
	ctx := context.Background()
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(userHomeDir(), "bitmark-node-data-test")
	if err != nil {
		return err
	}
	userPath = UserPath{
		BaseDir:        baseDir,
		NodeDBDir:      "db",
		MainnetDataDir: "data",
		TestnetDataDir: "data-test",
		MainnetLogDir:  "log",
		TestnetLogDir:  "log-test",
	}

	dockerPath = DockerPath{
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

	mock.Watcher = NodeWatcher{DockerClient: client, BackgroundContex: ctx,
		Repo:          "docker.io/bitmark/bitmark-node-test",
		ImageName:     "bitmark/bitmark-node-test",
		ContainerName: "bitmarkNodeTest",
		Postfix:       dockerPath.OldContainerPostfix}
	// Create Directory For Test
	mock.createDir()
	mock.Chain = "bitmark"
	// Create Environment Variable
	mock.Env = make(map[string]string)
	mock.Env["PUBLIC_IP"] = "0.0.0.0"
	mock.Env["NETWORK"] = mock.Chain
	mock.Env["DOCKER_HOST"] = "unix:///var/run/docker.sock"
	mock.Env["NODE_IMAGE"] = "bitmark/bitmark-node-test"
	mock.Env["NODE_NAME"] = "bitmarkNodeTest"
	mock.Env["USER_NODE_BASE_DIR"] = baseDir
	for k, v := range mock.Env {
		os.Setenv(k, v)
		fmt.Println("key:", k, " val:", v)
	}
	// Create sub directory  names
	return nil
}
func (mock *MockData) getWatcher() *NodeWatcher {
	return &mock.Watcher
}

func (mock *MockData) createDir() error {

	//Create file Dirs
	tempDir := mock.dockerPath.GetNodeDBPath()
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return err
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return err
	}
	tempDir = mock.dockerPath.GetDataPath(dockerPath.GetMainnet())
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return err
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return err
	}

	tempDir = mock.dockerPath.GetDataPath(dockerPath.GetTestnet())
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return err
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return err
	}

	tempDir = mock.dockerPath.GetLogPath(dockerPath.GetMainnet())
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return err
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return err
	}
	tempDir = mock.dockerPath.GetLogPath(dockerPath.GetTestnet())
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return err
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return err
	}
	return nil
}

// Utilies
func cleanContainers() error {
	// Clean Container
	watcher := mockData.getWatcher()
	containers, err := watcher.getContainersWithImage()
	if err != nil {
		return err
	}
	for _, container := range containers {
		fmt.Println("Container ID:", container.ID, " is going to be removed")
		watcher.forceRemoveContainer(container.ID)
	}
	containers, err = watcher.getContainersWithImage()
	if len(containers) > 0 {
		return errors.New("clean container error")
	}
	return nil
}
