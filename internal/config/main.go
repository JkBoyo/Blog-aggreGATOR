package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func Read() (Config, error) {
	confFP, err := getConfigFp()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(confFP)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var userConfig Config

	err = json.Unmarshal(data, &userConfig)
	if err != nil {
		return Config{}, err
	}

	return userConfig, nil
}

func getConfigFp() (string, error) {
	confFp, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confFp, configFileName), nil
}

func write(cfg Config) error {
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	fP, err := getConfigFp()
	if err != nil {
		return err
	}

	return os.WriteFile(fP, cfgJson, 0644)
}
