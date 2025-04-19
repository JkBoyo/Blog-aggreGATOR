package config

import ()

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
