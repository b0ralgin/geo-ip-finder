package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	DB       string                    `yaml:"db"`
	TTL      time.Duration             `yaml:"ttl"`
	Services map[string]GeoServicesCfg `yaml:"services"`
}

type GeoServicesCfg struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
	Limit uint16 `yaml:"limit"`
}

func LoadCfg(fileName string) (*Config, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
