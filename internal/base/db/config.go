package db

import (
	"fmt"
)

type Config struct {
	Port            uint16 `json:"port"              mapstructure:"port"              yaml:"port"`
	Host            string `json:"host"              mapstructure:"host"              yaml:"host"`
	User            string `json:"user"              mapstructure:"user"              yaml:"user"`
	Password        string `json:"password"          mapstructure:"password"          yaml:"password"`
	Schema          string `json:"schema"            mapstructure:"schema"            yaml:"schema"`
	Charset         string `json:"charset"           mapstructure:"charset"           yaml:"charset"`
	SslMode         string `json:"ssl_mode"          mapstructure:"ssl_mode"          yaml:"ssl_mode"`
	LogSQL          bool   `json:"log_sql"           mapstructure:"log_sql"           yaml:"log_sql"`
	MaxIdleConns    int    `json:"mas_idle_conns"    mapstructure:"max_idle_conns"    yaml:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns"    mapstructure:"max_open_conns"    yaml:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

func (c *Config) Address() string {
	conn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		c.User,
		c.Password,
		c.Schema,
		c.Host,
		c.Port,
		c.SslMode,
	)

	return conn
}
