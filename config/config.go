package config

import (
	"sync"
)

var once sync.Once
var instance *Config

type Config struct {
	Qbittorrent struct {
		URL      string
		Username string
		Password string
	}
	Utorrent struct {
		URL string
	}
	Nzbget struct {
		URL string
	}
	Port int
	Mode string
}

func Instance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func (c *Config) Validate() error {
	// Add validations for your configuration

	return nil
}
