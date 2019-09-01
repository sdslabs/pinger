package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const configPath = "config.yml"

// Config for `config.yml`
type Config struct {
}

// Parse takes the path of config file and uses a *Config to store data
func (c *Config) Parse(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

// GetConfig is shorthand for getting config from `config.yml`
func GetConfig() (Config, error) {
	c := &Config{}
	err := c.Parse(configPath)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}
