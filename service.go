package main

import (
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	log "github.com/google/logger"
)

const (
	containerStopWaitTime = 15 * time.Second
	pullImageInterval     = 20 * time.Second // Change to a longer Interval after finishing develop
	recoverWaitTime       = 20 * time.Second
)

// Database to remove
const (
	nodeDataDirMainnet  = "/.config/bitmark-node/bitmarkd/bitmark/data"
	nodeDataDirTestnet  = "/.config/bitmark-node/bitmarkd/testing/data"
	blockLevelDB        = "bitmark-blocks.leveldb"
	indexLevelDB        = "bitmark-index.leveldb"
	oldCotnainerPostfix = ".old"
	oldDBPostfix        = ".old"
)

// StartMonitor  Monitor process
func StartMonitor(watcher NodeWatcher) error {
	log.Info("Monitoring Process Start")
	defer func() {
		if r := recover(); r != nil {
			log.Error(ErrorStartMonitorService.Error() + ": waiting for restart service")
			time.Sleep(recoverWaitTime)
			go StartMonitor(watcher)
		}
	}()
	for {
		updated := make(chan bool)
		go imageUpdateRoutine(&watcher, updated)
		<-updated
		createConf, err := handleExistingContainer(watcher)
		if err != nil {
			log.Error(ErrCombind(ErrorHandleExistingContainer, err))
			continue
		}

		var newContainer container.ContainerCreateCreatedBody
		if createConf != nil { // err == nil and createConf == nil => container does not exist
			newContainer, err = watcher.createContainer(*createConf)
			if err != nil {
				log.Error(ErrCombind(ErrorContainerCreate, err))
				continue
			}
		} else {
			log.Info("Creating a brand new container")
			newContainerConfig, err := getDefaultConfig(&watcher)
			if err != nil {
				log.Error(ErrCombind(ErrorConfigCreateNew, err))
				continue
			}
			newContainer, err = watcher.DockerClient.ContainerCreate(watcher.BackgroundContex, newContainerConfig.Config,
				newContainerConfig.HostConfig, nil, watcher.ContainerName)
		}
		err = renameBitmarkDBs()
		if err != nil {
			log.Error(ErrCombind(ErrorRenameDB, err))
		}
		err = watcher.startContainer(newContainer.ID)
		if err != nil {
			log.Error(ErrCombind(ErrorContainerStart, err))
			err = recoverBitmarkDBs()
			if err != nil {
				log.Error(ErrCombind(ErrorRecoverDB, err))
			}
			continue
		}
		log.Info("Start container successfully")
	}
	return nil
}

// imageUpdateRoutine check image periodically
func imageUpdateRoutine(w *NodeWatcher, updateStatus chan bool) {
	ticker := time.NewTicker(pullImageInterval)
	defer func() {
		ticker.Stop()
		close(updateStatus)
	}()
	// For the first time
	newImage, err := w.pullImage()
	if err != nil {
		log.Info(ErrCombind(ErrorImagePull, err).Error())
	}
	if newImage {
		log.Info("imageUpdateRoutine update a new image")
		updateStatus <- true
	}
	// End of the first time ---- can be delete later
	for { // start  periodically check routine
		select {
		case <-ticker.C:
			newImage, err := w.pullImage()
			if err != nil {
				log.Info(ErrCombind(ErrorImagePull, err))
				continue
			}
			if newImage {
				log.Info("imageUpdateRoutine update a new image")
				updateStatus <- true
				break
			}
			log.Info("no new image found")
		}
	}
}

// handleExistingContainer handle existing container and return old container config for recreating a new container
func handleExistingContainer(watcher NodeWatcher) (*CreateConfig, error) {
	nodeContainers, err := watcher.getContainersWithImage()
	if err != nil { //not found is not an error
		log.Error(ErrCombind(ErrorGetContainerWithImage, err))
		return nil, nil
	}
	if len(nodeContainers) != 0 {
		nameContainer := watcher.getNamedContainer(nodeContainers)
		if nameContainer == nil { //not found is not an error
			log.Warning(ErrorNamedContainerNotFound.Error())
			return nil, nil
		}
		jsonConfig, err := watcher.DockerClient.ContainerInspect(watcher.BackgroundContex, nameContainer.ID)

		if err != nil { //inspect fail is an error because we can not do anything about existing error
			return nil, err
		}

		newConfig := container.Config{
			Image:        watcher.ImageName,
			ExposedPorts: jsonConfig.Config.ExposedPorts,
			Env:          jsonConfig.Config.Env,
			Volumes:      jsonConfig.Config.Volumes,
			Cmd:          jsonConfig.Config.Cmd,
		}

		newNetworkConf := network.NetworkingConfig{
			EndpointsConfig: jsonConfig.NetworkSettings.Networks,
		}

		namedContainers := append([]types.Container{}, *nameContainer)
		err = watcher.stopContainers(namedContainers, containerStopWaitTime)

		if err != nil {
			return nil, err
		}

		oldContainers, err := watcher.getOldContainer()
		if err == nil && oldContainers != nil {
			watcher.forceRemoveContainer(oldContainers.ID)
		}
		err = watcher.renameContainer(nameContainer)
		if err != nil {
			return nil, err
		}
		return &CreateConfig{Config: &newConfig, HostConfig: jsonConfig.HostConfig, NetworkingConfig: &newNetworkConf}, err
	}
	// no container
	return nil, nil
}

func getDefaultConfig(watcher *NodeWatcher) (*CreateConfig, error) {
	config := CreateConfig{}

	baseDir, err := builDefaultVolumSrcBaseDir(watcher)
	log.Info("baseDir:", baseDir)
	if err != nil {
		return nil, err
	}

	baseTargetDir := "/.config/bitmark-node"
	publicIP := os.Getenv("PUBLIC_IP")
	chain := os.Getenv("NETWORK")
	if len(publicIP) == 0 {
		publicIP = "127.0.0.1"
	}
	if len(chain) == 0 {
		chain = "BITMARK"
	}
	additionEnv := append([]string{}, "PUBLIC_IP="+publicIP, "NETWORK="+chain)
	exposePorts := nat.PortMap{
		"2136/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "2136",
			},
		},
		"2130/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "2130",
			},
		},
		"2131/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "2131",
			},
		},
		"9980/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "9980",
			},
		},
	}

	hconfig := container.HostConfig{
		NetworkMode:  "default",
		PortBindings: exposePorts,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: baseDir + "/db",
				Target: baseTargetDir + "/db",
			},
			{
				Type:   mount.TypeBind,
				Source: baseDir + "/data",
				Target: baseTargetDir + "/bitmarkd/bitmark/data",
			},
			{
				Type:   mount.TypeBind,
				Source: baseDir + "/data-test",
				Target: baseTargetDir + "/bitmarkd/testing/data",
			},
			{
				Type:   mount.TypeBind,
				Source: baseDir + "/log",
				Target: baseTargetDir + "/bitmarkd/bitmark/log",
			},
			{
				Type:   mount.TypeBind,
				Source: baseDir + "/log-test",
				Target: baseTargetDir + "/bitmarkd/testing/log",
			},
		},
	}
	config.HostConfig = &hconfig
	portmap := nat.PortSet{
		"2136/tcp": struct{}{},
		"2130/tcp": struct{}{},
		"2131/tcp": struct{}{},
		"9980/tcp": struct{}{},
	}
	config.Config = &container.Config{
		Image:        watcher.ImageName,
		Env:          additionEnv,
		ExposedPorts: portmap,
	}

	return &config, nil
}

func renameDB(src, dest string) (err error) {
	_, err = os.Stat(src)
	if nil == err {
		if err = os.Rename(src, dest); err != nil {
			return err
		}
	}
	return err
}

func renameBitmarkDBs() (finalErr error) {
	mainnetBlockDB := nodeDataDirMainnet + "/" + blockLevelDB
	if err := renameDB(mainnetBlockDB, mainnetBlockDB+oldDBPostfix); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}

	mainnetIndexDB := nodeDataDirMainnet + "/" + indexLevelDB
	if err := renameDB(mainnetIndexDB, mainnetIndexDB+oldDBPostfix); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}

	testnetBlockDB := nodeDataDirTestnet + "/" + blockLevelDB
	if err := renameDB(testnetBlockDB, testnetBlockDB+oldDBPostfix); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}

	testnetIndexDB := nodeDataDirTestnet + "/" + indexLevelDB
	if err := renameDB(testnetIndexDB, testnetIndexDB+oldDBPostfix); err != nil {
		finalErr = ErrCombind(finalErr, err)
	}
	return finalErr
}

func builDefaultVolumSrcBaseDir(watcher *NodeWatcher) (string, error) {
	homeDir := os.Getenv("USER_NODE_BASE_DIR")
	if 0 == len(homeDir) {
		return "", ErrorUserNodeDirEnv
	}
	return homeDir, nil
}

func recoverBitmarkDBs() (finalErr error) {
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
