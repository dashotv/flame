package app

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/pkg/errors"
)

var cfg *Config

func setupConfig() (err error) {
	cfg = &Config{}
	if err := env.Parse(cfg); err != nil {
		return errors.Wrap(err, "failed to parse environment variables")
	}

	// fmt.Println("Connections:")
	// for k, v := range cfg.Connections {
	// 	fmt.Printf("  %15s: %+v\n", k, v)
	// }

	if err := cfg.Validate(); err != nil {
		fmt.Printf("validation failed: %+v\n", cfg)
		return errors.Wrap(err, "failed to validate config")
	}

	return nil
}

type Config struct {
	Port                int    `env:"PORT"`
	Mode                string `env:"MODE"`
	Logger              string `env:"LOGGER"`
	Cron                bool   `env:"CRON"`
	QbittorrentURL      string `env:"QBITTORRENT_URL"`
	QbittorrentUsername string `env:"QBITTORRENT_USERNAME"`
	QbittorrentPassword string `env:"QBITTORRENT_PASSWORD"`
	NzbgetURL           string `env:"NZBGET_URL"`
	NatsURL             string `env:"NATS_URL"`
	RedisHost           string `env:"REDIS_HOST"`
	RedisPort           string `env:"REDIS_PORT"`
	RedisDatabase       int    `env:"REDIS_DATABASE"`
}

type Connection struct {
	URI        string `yaml:"uri,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Collection string `yaml:"collection,omitempty"`
}

func (c *Connection) UnmarshalText(text []byte) error {
	vals := strings.Split(string(text), ",")
	c.URI = vals[0]
	c.Database = vals[1]
	c.Collection = vals[2]
	return nil
}

type ConnectionSet map[string]*Connection

func (c *ConnectionSet) UnmarshalText(text []byte) error {
	*c = make(map[string]*Connection)
	for _, conn := range strings.Split(string(text), ";") {
		kv := strings.Split(conn, "=")
		vals := strings.Split(kv[1]+",,", ",")
		(*c)[kv[0]] = &Connection{
			URI:        vals[0],
			Database:   vals[1],
			Collection: vals[2],
		}
	}
	return nil
}

func (c *Config) RedisURL() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Config) Validate() error {
	// No Connections
	// if err := c.validateDefaultConnection(); err != nil {
	// 	return err
	// }
	// TODO: add more validations?
	return nil
}
