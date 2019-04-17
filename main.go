package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	log "github.com/google/logger"
	"github.com/urfave/cli"
)

const (
	updaterVersion   string = "1.0.0"
	appName          string = "bitmark-node-updater"
	dockerAPIVersion string = "1.24"
	logPath          string = "bitmark-node-watcher.log"
)

var (
	userPath   UserPath
	dockerPath DockerPath
)

func main() {
	// assign it to the standard logger
	app := cli.NewApp()
	app.Name = appName
	app.Version = updaterVersion
	app.Usage = "Automatically update running bitmark-node container"
	app.Before = before
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, H",
			Usage:  "daemon socket to connect to",
			Value:  "unix:///var/run/docker.sock",
			EnvVar: "DOCKER_HOST",
		},

		cli.StringFlag{
			Name:  "image, i",
			Usage: "image name to pull",
			Value: "bitmark/bitmark-node",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "container name to create",
			Value: "bitmarkNode",
		},
		cli.BoolFlag{
			Name:  "verbose, vb",
			Usage: "verbose of log",
		},
	}

	app.Action = func(c *cli.Context) error {
		logfile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}
		verbose := c.GlobalBool("verbose")
		log.Init("bitmark-node-updater-log", verbose, false, logfile)
		log.Info("create log file in bitmark-node-updater-log")
		defer logfile.Close()

		ctx := context.Background()
		client, err := client.NewEnvClient()
		if err != nil {
			log.Error(ErrorGetAPIFail)
			return err
		}
		// Create a Docker API Client and current Context
		dockerImage := c.GlobalString("image")
		dockerRepo := "docker.io/" + dockerImage
		containerName := c.GlobalString("name")
		baseDir, err := builDefaultVolumSrcBaseDir()
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

		watcher := NodeWatcher{DockerClient: client, BackgroundContex: ctx,
			Repo: dockerRepo, ImageName: dockerImage, ContainerName: containerName, Postfix: dockerPath.OldContainerPostfix}
		err = StartMonitor(watcher)
		if err != nil {
			log.Errorf(ErrorStartMonitorService.Error(), " image name:", watcher.ImageName)
			return err
		}
		log.Info("Start Monitor host:", c.GlobalString("host"), "image:", c.GlobalString("image"))
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func before(c *cli.Context) error {
	// configure environment vars for client
	err := envConfig(c)
	if err != nil {
		log.Info("envConfig Error", err)
		return err
	}
	return nil
}

// envConfig translates the command-line options into environment variables
// that will initialize the api client
func envConfig(c *cli.Context) error {
	if err := setEnvOptStr("DOCKER_HOST", c.GlobalString("host")); err != nil {
		return err
	}
	if err := setEnvOptStr("DOCKER_API_VERSION", dockerAPIVersion); err != nil {
		return err
	}
	return nil
}

func setEnvOptStr(env string, opt string) error {
	if opt != "" && opt != os.Getenv(env) {
		err := os.Setenv(env, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func setEnvOptBool(env string, opt bool) error {
	if opt == true {
		return setEnvOptStr(env, "1")
	}
	return nil
}
