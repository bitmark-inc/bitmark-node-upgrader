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
	// Time Interval
	containerStopWaitTime = 15 * time.Second
	pullImageInterval     = 120 * time.Second // Change to a longer Interval after finishing develop
	recoverWaitTime       = 20 * time.Second
	runAPIDelayTime       = 5 * time.Second
	// UpdateDB Endpoint
	updateDBEndPoint = "https://0u08da25ba.execute-api.ap-northeast-1.amazonaws.com/v1/chaindata"
)

// StartMonitor  Monitor process
func StartMonitor(watcher NodeWatcher) error {
	log.Info("Monitoring Process Start")
	defer func() {
		log.Warning("StartMonitor Exit")
		if r := recover(); r != nil {
			log.Error(ErrorStartMonitorService.Error() + ": waiting for restart service")
			time.Sleep(recoverWaitTime)
			go StartMonitor(watcher)
		}
	}()
	updateDBConf, err := makeUpdateDBConfig()
	if err != nil {
		log.Error(err)
		return err
	}
	// Force check CheckUpdateDB for existing database
	err = firstTimeUpdateDB(&watcher)
	if err != nil {
		log.Error(ErrCombind(ErrorFirstUpdateDB, err))
	}
	for {
		updated := make(chan bool)
		go imageUpdateCheckRoutine(&watcher, updated)
		<-updated
		log.Info("Start to Update New Image")
		createConf, err := handleExistingContainer(watcher)
		if err != nil {
			log.Error(ErrCombind(ErrorHandleExistingContainer, err))
			continue
		}

		var newContainer container.ContainerCreateCreatedBody
		if createConf != nil { // err == nil and createConf == nil => container does not exist
			log.Info("Creating a container from existing container")
			newContainer, err = watcher.createContainer(*createConf)
			if err != nil {
				log.Error(ErrCombind(ErrorContainerCreate, err))
				continue
			}
		} else {
			log.Info("Creating a brand new container")
			var newContainerConfig *CreateContainerConfig
			newContainerConfig, err = getDefaultConfig(&watcher)
			if err != nil {
				log.Error(ErrCombind(ErrorConfigCreateNew, err))
				continue
			}
			newContainer, err = watcher.createContainer(*newContainerConfig)
			if err != nil {
				log.Error(ErrCombind(ErrorContainerCreate, err))
				continue
			}
		}
		dbUpdater, err := SetDBUpdaterReady(updateDBConf)
		if err != nil {
			log.Error(ErrCombind(ErrorSetDBUpdaterReady, err))
		}
		err = dbUpdater.(*DBUpdaterHTTPS).UpdateToLatestDB()
		if err != nil {
			log.Error(ErrCombind(ErrorUpdateToLatestDB, err))
		}
		err = watcher.startContainer(newContainer.ID)
		if err != nil {
			log.Error(ErrCombind(ErrorContainerStart, err))
			err = recoverBitmarkdDB()
			if err != nil {
				log.Error(ErrCombind(ErrorRecoverDB, err))
			}
			continue
		} else { // No container error start services
			go runBitmarkdRecorderd()
		}
		log.Info("Start container successfully")
		time.Sleep(pullImageInterval) // This run start immediately
	}
}

// imageUpdateRoutine check image periodically
func imageUpdateCheckRoutine(w *NodeWatcher, updateStatus chan bool) {
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
		log.Info("imageUpdateCheckRoutine found a new image")
		updateStatus <- true
		return
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
				log.Info("imageUpdateCheckRoutine found a new image")
				updateStatus <- true
				break
			}
			// TODO: remove log
			log.Info("no new image found")
		}
	}
}

func firstTimeUpdateDB(watcher *NodeWatcher) (err error) {
	nodeContainers, err := watcher.getContainersWithImage()
	log.Warning("start firstTimeUpdateDB check...")
	if err != nil { //not found is not an error
		log.Error(ErrCombind(ErrorGetContainerWithImage, err))
		return nil
	}
	if len(nodeContainers) != 0 { // Container exist
		log.Warning("container found start to upgrade db process ...")
		nameContainer := watcher.getNamedContainer(nodeContainers)
		if nameContainer == nil { //not found is not an error
			log.Warning(ErrorNamedContainerNotFound.Error())
			return nil
		}
		namedContainers := append([]types.Container{}, *nameContainer)
		err = watcher.stopContainers(namedContainers, containerStopWaitTime)
		if err != nil { //inspect fail is an error because we can not do anything about existing error
			return err
		}
		updateDBConf, err := makeUpdateDBConfig()
		dbUpdater, err := SetDBUpdaterReady(updateDBConf)
		if err != nil {
			log.Error(ErrCombind(ErrorSetDBUpdaterReady, err))
		}
		err = dbUpdater.(*DBUpdaterHTTPS).UpdateToLatestDB()
		if err != nil {
			log.Error(ErrCombind(ErrorUpdateToLatestDB, err))
		}
		err = watcher.startContainer(nameContainer.ID)
		if err != nil {
			log.Error(ErrCombind(ErrorContainerStart, err))
			err = recoverBitmarkdDB()
			if err != nil {
				log.Error(ErrCombind(ErrorRecoverDB, err))
			}
		} else { // Start Container Successfully, start bitmarkd and recorderd
			runBitmarkdRecorderd()
		}
	}
	return nil
}

// handleExistingContainer handle existing container and return old container config for recreating a new container
func handleExistingContainer(watcher NodeWatcher) (*CreateContainerConfig, error) {
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
		return &CreateContainerConfig{Config: &newConfig, HostConfig: jsonConfig.HostConfig, NetworkingConfig: &newNetworkConf}, err
	}
	// no container
	return nil, nil
}

func getDefaultConfig(watcher *NodeWatcher) (*CreateContainerConfig, error) {
	config := CreateContainerConfig{}

	publicIP := os.Getenv("PUBLIC_IP")
	chain := os.Getenv("NETWORK")
	if len(publicIP) == 0 {
		publicIP = "127.0.0.1"
	}
	if len(chain) == 0 {
		chain = "bitmark"
	}
	log.Info("PUBLIC_IP="+publicIP, "NETWORK="+chain)
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
				Source: userPath.GetNodeDBPath(),
				Target: dockerPath.GetNodeDBPath(),
			},
			{
				Type:   mount.TypeBind,
				Source: userPath.GetDataPath(userPath.GetMainnet()),
				Target: dockerPath.GetDataPath(userPath.GetMainnet()),
			},
			{
				Type:   mount.TypeBind,
				Source: userPath.GetDataPath(userPath.GetTestnet()),
				Target: dockerPath.GetDataPath(userPath.GetTestnet()),
			},
			{
				Type:   mount.TypeBind,
				Source: userPath.GetLogPath(userPath.GetMainnet()),
				Target: dockerPath.GetLogPath(dockerPath.GetMainnet()),
			},
			{
				Type:   mount.TypeBind,
				Source: userPath.GetLogPath(userPath.GetTestnet()),
				Target: dockerPath.GetLogPath(dockerPath.GetTestnet()),
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

func makeUpdateDBConfig() (DBUpdaterHTTPSConfig, error) {
	config := DBUpdaterHTTPSConfig{
		APIEndpoint:        updateDBEndPoint,
		CurrentDBPath:      dockerPath.GetBlockDBPath(userPath.GetMainnet()),
		ZipSourcePath:      dockerPath.GetUpdateDBZipFilePath(dockerPath.GetMainnet()),
		ZipDestinationPath: dockerPath.GetDataPath(dockerPath.GetMainnet()),
	}
	return config, nil
}

func runBitmarkdRecorderd() {
	time.Sleep(runAPIDelayTime)
	NodeAPI("", bitmarkdStart)
	NodeAPI("", recorderdStart)
}
