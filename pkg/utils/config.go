package utils

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configPath = "config.yml"

// Config is the data from `config.yml`
var Config ConfigFile

// OauthProviderConfig provides configuration settings for an OAuth Provider.
type OauthProviderConfig struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURL  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes"`
}

type oauth = map[string]OauthProviderConfig

type database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  bool   `yaml:"ssl_mode"`
}

type application struct {
	Secret   string   `yaml:"secret"`
	Oauth    oauth    `yaml:"oauth"`
	Database database `yaml:"database"`
}

type central struct{}

// ConfigFile for `config.yml`
type ConfigFile struct {
	Application application `yaml:"application"`
	Central     central     `yaml:"central"`
}

// Parse takes the path of config file and uses a *Config to store data
func (c *ConfigFile) Parse(path string) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

// GetConfig is shorthand for getting config from `config.yml`
func GetConfig() (ConfigFile, error) {
	c := &ConfigFile{}
	err := c.Parse(configPath)
	if err != nil {
		return ConfigFile{}, err
	}
	return *c, nil
}
