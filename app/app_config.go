package app

import (
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/pkg/errors"
)

func init() {
	initializers = append(initializers, setupConfig)
}

func setupConfig(app *Application) error {
	app.Config = &Config{}
	if err := env.Parse(app.Config); err != nil {
		return errors.Wrap(err, "parsing config")
	}

	if err := app.Config.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate config")
	}

	return nil
}

type Config struct {
	Mode                string `env:"MODE" envDefault:"dev"`
	Logger              string `env:"LOGGER" envDefault:"dev"`
	Port                int    `env:"PORT" envDefault:"10080"`
	QbittorrentURL      string `env:"QBITTORRENT_URL"`
	QbittorrentUsername string `env:"QBITTORRENT_USERNAME"`
	QbittorrentPassword string `env:"QBITTORRENT_PASSWORD"`
	NzbgetURL           string `env:"NZBGET_URL"`
	RedisHost           string `env:"REDIS_HOST"`
	RedisPort           string `env:"REDIS_PORT"`
	RedisDatabase       int    `env:"REDIS_DATABASE"`

	//golem:template:app/config_partial_struct
	// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
	// Models (Database)
	Connections ConnectionSet `env:"CONNECTIONS"`

	// Cache
	RedisAddress string `env:"REDIS_ADDRESS"`

	// Router Auth
	Auth           bool   `env:"AUTH" envDefault:"false"`
	ClerkSecretKey string `env:"CLERK_SECRET_KEY"`

	// Events
	NatsURL string `env:"NATS_URL"`

	// Workers
	MinionConcurrency int    `env:"MINION_CONCURRENCY" envDefault:"10"`
	MinionDebug       bool   `env:"MINION_DEBUG" envDefault:"false"`
	MinionBufferSize  int    `env:"MINION_BUFFER_SIZE" envDefault:"100"`
	MinionURI         string `env:"MINION_URI"`
	MinionDatabase    string `env:"MINION_DATABASE"`
	MinionCollection  string `env:"MINION_COLLECTION"`

	//golem:template:app/config_partial_struct

}

func (c *Config) Validate() error {
	list := []func() error{
		c.validateMode,
		c.validateLogger,
		//golem:template:app/config_partial_validate
		// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
		c.validateDefaultConnection,

		//golem:template:app/config_partial_validate

	}

	for _, fn := range list {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) validateMode() error {
	switch c.Mode {
	case "dev", "release":
		return nil
	default:
		return errors.New("invalid mode (must be 'dev' or 'release')")
	}
}

func (c *Config) validateLogger() error {
	switch c.Logger {
	case "dev", "release":
		return nil
	default:
		return errors.New("invalid logger (must be 'dev' or 'release')")
	}
}

//golem:template:app/config_partial_connection
// DO NOT EDIT. This section is managed by github.com/dashotv/golem.

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

func (c *Config) ConnectionFor(name string) (*Connection, error) {
	def, ok := c.Connections["default"]
	if !ok {
		return nil, errors.Errorf("connection for %s: no default connection found", name)
	}

	conn, ok := c.Connections[name]
	if !ok {
		return nil, errors.Errorf("no connection named '%s'", name)
	}

	if conn.URI == "" {
		conn.URI = def.URI
	}
	if conn.Database == "" {
		conn.Database = def.Database
	}
	if conn.Collection == "" {
		conn.Collection = def.Collection
	}

	return conn, nil
}

//golem:template:app/config_partial_connection
