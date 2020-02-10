package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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
	secret   string   `yaml:"secret"`
	Oauth    oauth    `yaml:"oauth"`
	Database database `yaml:"database"`
}

type central struct{}

// StatusConfig for `config.yml`
type StatusConfig struct {
	Application application `yaml:"application"`
	Central     central     `yaml:"central"`
}

// Secret returns the secret key to encrypt tokens.
func (c *StatusConfig) Secret() []byte {
	return []byte(c.Application.secret)
}

// Parse takes the path of config file and uses a *Config to store data
func (c *StatusConfig) Parse(path string) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

// GetStatusConfig is shorthand for getting config from `config.yml`
func GetStatusConfig(configPath string) (StatusConfig, error) {
	c := &StatusConfig{}
	err := c.Parse(configPath)
	if err != nil {
		return StatusConfig{}, err
	}
	return *c, nil
}
