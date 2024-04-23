package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFilename = ".pmrc"

type Config struct {
	DefaultPackageManager string `json:"defaultPackageManager"`
}

var defaultConfig = Config{
	DefaultPackageManager: "npm",
}

func ReadConfig() Config {
	homeDir, _ := os.UserHomeDir()
	path, err := filepath.Abs(filepath.Join(homeDir, configFilename))
	if err != nil {
		return defaultConfig
	}

	file, err := os.Open(path)
	if err != nil {
		return defaultConfig
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return defaultConfig
	}

	if config.DefaultPackageManager == "" {
		config.DefaultPackageManager = defaultConfig.DefaultPackageManager
	}

	return config
}
