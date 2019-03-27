package main

import (
	"errors"
	"fmt"
)

var ( // Error variable
	// Directory Error
	ErrorUserNodeDirEnv = errors.New("User input node base directory not found")
	ErrorRenameDB       = errors.New("rename db failed")
	ErrorRecoverDB      = errors.New("rename db failed")
	// Process Error
	ErrorGetAPIFail              = errors.New("Get Docker API failed")
	ErrorStartMonitorService     = errors.New("StartMonitor failed")
	ErrorListenContainer         = errors.New("Listen Container failed")
	ErrorHandleExistingContainer = errors.New("Handle existing container failed")
	ErrorImageUpdateRoutine      = errors.New("Image update routine failed")
	// Container Errors
	ErrorContainerCreate        = errors.New("Container create failed")
	ErrorContainerStart         = errors.New("Container start  failed")
	ErrorContainerStop          = errors.New("Container stop failed")
	ErrorConfigCreateNew        = errors.New("Create a new Config error")
	ErrorNamedContainerNotFound = errors.New("Named container is not found")
	// Image Errors
	ErrorImagePull             = errors.New("Image pull failed")
	ErrorGetContainerWithImage = errors.New("Get container with image failed")

	// NodeWatcher Errors
	ErrorCreateWatcher = errors.New("Create NodeWatcher failed")
)

// Combine all Errors together
func ErrCombind(errs ...error) (finalErr error) {
	for _, err := range errs {
		if err != nil {
			if finalErr != nil {
				finalErr = fmt.Errorf("%s-%s", finalErr, err)
			} else {
				finalErr = fmt.Errorf("%s", err)
			}
		}
	}
	return finalErr
}
