package main

import "path/filepath"

const (
	// networks
	mainnet = "bitmark"
	testnet = "testnet"
)

// GetMainnet return mainnet string to present network
func (u *UserPath) GetMainnet() string {
	return mainnet
}

// GetTestnet return testnet string to present network
func (u *UserPath) GetTestnet() string {
	return testnet
}

// GetNodeDBPath return full data path accordint to network
func (u *UserPath) GetNodeDBPath() string {
	return filepath.Join(u.BaseDir, u.NodeDBDir)
}

// GetDataPath return full data path accordint to network
func (u *UserPath) GetDataPath(network string) (path string) {
	if mainnet == network {
		path = filepath.Join(u.BaseDir, u.MainnetDataDir)
	} else {
		path = filepath.Join(u.BaseDir, u.TestnetDataDir)
	}
	return path
}

// GetLogPath return full log path accordint to network
func (u *UserPath) GetLogPath(network string) (path string) {
	if mainnet == network {
		path = filepath.Join(u.BaseDir, u.MainnetLogDir)
	} else {
		path = filepath.Join(u.BaseDir, u.TestnetLogDir)
	}
	return path
}

// GetMainnet return mainnet string to present network
func (d *DockerPath) GetMainnet() string {
	return mainnet
}

// GetTestnet return testnet string to present network
func (d *DockerPath) GetTestnet() string {
	return testnet
}

// GetNodeDBPath return full data path accordint to network
func (d *DockerPath) GetNodeDBPath() string {
	return filepath.Join(d.BaseDir, d.NodeDBDir)
}

// GetDataPath return full data path accordint to network
func (d *DockerPath) GetDataPath(network string) (path string) {
	if mainnet == network {
		path = filepath.Join(d.BaseDir, d.MainnetDataDir)
	} else {
		path = filepath.Join(d.BaseDir, d.TestnetDataDir)
	}
	return path
}

// GetLogPath return full log path accordint to network
func (d *DockerPath) GetLogPath(network string) (path string) {
	if mainnet == network {
		path = filepath.Join(d.BaseDir, d.MainnetLogDir)
	} else {
		path = filepath.Join(d.BaseDir, d.TestnetLogDir)
	}
	return path
}

// GetBlockDBPath take network name (mainnet or testnet) and return block leveldb path
func (d *DockerPath) GetBlockDBPath(network string) string {
	return filepath.Join(d.GetDataPath(network), d.BlockDBDirName)
}

// GetIndexDBPath take network name (mainnet or testnet) and return index leveldb path
func (d *DockerPath) GetIndexDBPath(network string) string {
	return filepath.Join(d.GetDataPath(network), d.IndexDBDirName)
}

// GetUpdateDBZipFilePath take network name (mainnet or testnet) and return index leveldb path
func (d *DockerPath) GetUpdateDBZipFilePath(network string) string {
	return filepath.Join(d.GetDataPath(network), d.UpdateDBZipName)
}

// AppendContainerPostfixToPath take network name (mainnet or testnet) and return index leveldb path
func (d *DockerPath) AppendContainerPostfixToPath(path string) string {
	return filepath.Join(path, d.OldContainerPostfix)
}

// AppendDBPostfixToPath take network name (mainnet or testnet) and return index leveldb path
func (d *DockerPath) AppendDBPostfixToPath(path string) string {
	return filepath.Join(path, d.OldDatabasePostfix)
}
