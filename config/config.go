package config

import (
	"errors"
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
	Port        int
	Mode        string
	Cron        bool
	Connections map[string]*Connection `yaml:"connections"`
}

type Connection struct {
	URI        string `yaml:"uri,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Collection string `yaml:"collection,omitempty"`
}

func Instance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func (c *Config) Validate() error {
	if err := c.validateDefaultConnection(); err != nil {
		return err
	}
	// TODO: add more validations?
	return nil
}

func (c *Config) validateDefaultConnection() error {
	if len(c.Connections) == 0 {
		return errors.New("you must specify a default connection")
	}

	var def *Connection
	for n, c := range c.Connections {
		if n == "default" || n == "Default" {
			def = c
			break
		}
	}

	if def == nil {
		return errors.New("no 'default' found in connections list")
	}
	if def.Database == "" {
		return errors.New("default connection must specify database")
	}
	if def.URI == "" {
		return errors.New("default connection must specify URI")
	}

	return nil
}
