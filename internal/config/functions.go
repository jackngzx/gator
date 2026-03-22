package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Cannot get User Home Directory")
	}
	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Cannot get config file path")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("Cannot read data from file")
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("Error decoding data from config file")
	}

	return cfg, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Cannot get config file path")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error opening the file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("Cannot encode JSON into file")
	}
	return nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return write(*cfg)
}
