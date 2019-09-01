package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const configPath = "config.yml"

type googleOauth struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURL  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes"`
}

type oauth struct {
	Google googleOauth `yaml:"google"`
}

// Config for `config.yml`
type Config struct {
	JWTSecret string `yaml:"jwt_secret"`
	Oauth     oauth  `yaml:"oauth"`
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
