package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server *Server `json:"server" mapstructure:"server" yaml:"server"`
		DB     *DB     `json:"db"     mapstructure:"db"     yaml:"db"`
		Key    *Key    `json:"key"    mapstructure:"key"    yaml:"key"`
		I18n   *I18n   `json:"i18n"   mapstructure:"i18n"   yaml:"i18n"`
	}

	Server struct {
		ShowStartBanner bool `json:"show_start_banner" mapstructure:"show_start_banner" yaml:"show_start_banner"`

		Addr string `json:"addr" mapstructure:"addr" yaml:"addr"`
	}

	DB struct {
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

	Key struct {
		Encryption      string `json:"encryption" mapstructure:"encryption" yaml:"encryption"`
		Hash            string `json:"hash"       mapstructure:"hash"       yaml:"hash"`
		EncryptionBytes []byte `json:"-"          mapstructure:"-"          yaml:"-"`
		HashBytes       []byte `json:"-"          mapstructure:"-"          yaml:"-"`
	}

	I18n struct {
		BundleDir string `json:"bundle_dir" mapstructure:"bundle_dir" yaml:"bundle_dir"`
	}
)

func LoadConfig() (*Config, error) {
	var configPath string

	env := os.Getenv(constants.AppEnv)
	if env == "" {
		env = constants.Dev
	}

	configPathFromEnv := os.Getenv(constants.ConfigPath)
	if configPathFromEnv != "" {
		configPath = configPathFromEnv
	} else {
		rootPath, err := getConfigRootPath()
		if err != nil {
			return nil, err
		}
		configPath = fmt.Sprintf("%s/config.development.yaml", rootPath)
	}

	fmt.Printf("Load config from %s\n", configPath)

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	secretEncryptionKeyBytes, err := base64.StdEncoding.DecodeString(cfg.Key.Encryption)
	if err != nil {
		return nil, fmt.Errorf("Could not decode email-encryption-key: %v", err)
	}
	cfg.Key.EncryptionBytes = secretEncryptionKeyBytes

	hashingKeyBytes, err := base64.StdEncoding.DecodeString(cfg.Key.Hash)
	if err != nil {
		return nil, fmt.Errorf("Could not decode email-hash-key: %v", err)
	}
	cfg.Key.HashBytes = hashingKeyBytes

	return cfg, nil
}

func (c *DB) Address() string {
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

func getConfigRootPath() (string, error) {
	// Get the current working directory
	// Getwd gives us the current working directory that we are running our app with `go run` command. in goland we can specify this working directory for the project
	// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// fmt.Printf("Current working directory is: %s\n", wd)

	// Get the absolute path of the executed project directory
	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "config")

	return configPath, nil
}
