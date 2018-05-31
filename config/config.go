package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Db       string                    `yaml:"db"`
	TTL      time.Duration             `yaml:"ttl"`
	Services map[string]GeoServicesCfg `yaml:"services"`
}

type GeoServicesCfg struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
	Limit int64  `yaml:"limit"`
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
