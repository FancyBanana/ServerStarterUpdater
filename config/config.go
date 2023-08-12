package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ModpackSlug string `yaml:"modpackSlug"`
	ModpackId   int    `yaml:"modpackId"`
	ApiKey      string `yaml:"apiKey"`
}

func ReadConfig(configPath string) (*Config, error) {

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	conf := Config{}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return nil, err
	}

	if conf.ModpackId != 0 && conf.ModpackSlug != "" {
		return nil, errors.New("Cannot have ModpackSlug and ModpackId both defined due to ambiguity")
	}

	if conf.ModpackId == 0 && conf.ModpackSlug == "" {
		return nil, errors.New("Must have ModpackSlug or ModpackId defined")
	}

	if conf.ApiKey == "" {
		return nil, errors.New("Must have apiKey defined")
	}

	return &conf, nil
}
