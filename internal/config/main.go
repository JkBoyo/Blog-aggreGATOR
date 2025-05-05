package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetDB(db string) error {
	c.DbUrl = db
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	confFP, err := getConfigFp()
	if err != nil {
		return Config{}, err
	}
	defaultConf := Config{}
	_, err = os.Stat(confFP)
	if os.IsNotExist(err) {
		err = write(defaultConf)
		if err != nil {
			return Config{}, err
		}
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
