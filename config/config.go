package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/spf13/viper"
)

var AWASM_URL string

type (
	Config struct {
		IngressURL string  `json:"ingressURL,omitempty" mapstructure:"ingress_url" yaml:"ingress_url,omitempty"`
		Server     *Server `json:"server,omitempty"     mapstructure:"server"      yaml:"server,omitempty"`
		DB         *DB     `json:"db,omitempty"         mapstructure:"db"          yaml:"db,omitempty"`
		Redis      *Redis  `json:"redis,omitempty"      mapstructure:"redis"       yaml:"redis,omitempty"`
		JWT        *JWT    `json:"jwt,omitempty"        mapstructure:"jwt"         yaml:"jwt,omitempty"`
		Key        *Key    `json:"key,omitempty"        mapstructure:"key"         yaml:"key,omitempty"`
		I18n       *I18n   `json:"i18n,omitempty"       mapstructure:"i18n"        yaml:"i18n,omitempty"`
	}

	Server struct {
		Addr            string `json:"addr,omitempty"  mapstructure:"addr"              yaml:"addr,omitempty"`
		ShowStartBanner bool   `json:"showStartBanner" mapstructure:"show_start_banner" yaml:"show_start_banner"`
	}

	DB struct {
		Port            uint16 `json:"port,omitempty"     mapstructure:"port"              yaml:"port,omitempty"`
		Host            string `json:"host,omitempty"     mapstructure:"host"              yaml:"host,omitempty"`
		User            string `json:"user,omitempty"     mapstructure:"user"              yaml:"user,omitempty"`
		Password        string `json:"password,omitempty" mapstructure:"password"          yaml:"password,omitempty"`
		Schema          string `json:"schema,omitempty"   mapstructure:"schema"            yaml:"schema,omitempty"`
		Charset         string `json:"charset"            mapstructure:"charset"           yaml:"charset"`
		SslMode         string `json:"sslMode"            mapstructure:"ssl_mode"          yaml:"ssl_mode"`
		LogSQL          bool   `json:"logSql"             mapstructure:"log_sql"           yaml:"log_sql"`
		MaxIdleConns    int    `json:"maxIdleConns"       mapstructure:"max_idle_conns"    yaml:"max_idle_conns"`
		MaxOpenConns    int    `json:"maxOpenConns"       mapstructure:"max_open_conns"    yaml:"max_open_conns"`
		ConnMaxLifetime int    `json:"connMaxLifetime"    mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
		MigrationDir    string `json:"migrationDir"       mapstructure:"migration_dir"     yaml:"migration_dir"`
	}

	Redis struct {
		Host       string `json:"host,omitempty"       mapstructure:"host"        yaml:"host,omitempty"`
		Port       int    `json:"port,omitempty"       mapstructure:"port"        yaml:"port,omitempty"`
		RequireTLS bool   `json:"requireTLS,omitempty" mapstructure:"require_tls" yaml:"require_tls,omitempty"`
		Password   string `json:"password"             mapstructure:"password"    yaml:"password"`
		DB         int    `json:"db,omitempty"         mapstructure:"db"          yaml:"db,omitempty"`
		PoolSize   int    `json:"poolSize"             mapstructure:"pool_size"   yaml:"pool_size"`
	}

	JWT struct {
		Secret string `json:"secret,omitempty" mapstructure:"secret" yaml:"secret,omitempty"`
		Exp    int    `json:"exp,omitempty"    mapstructure:"exp"    yaml:"exp,omitempty"`
	}

	Key struct {
		Encryption      string `json:"encryption,omitempty" mapstructure:"encryption"       yaml:"encryption,omitempty"`
		Hash            string `json:"hash,omitempty"       mapstructure:"hash"             yaml:"hash,omitempty"`
		EncryptionBytes []byte `json:"-"                    mapstructure:"encryption_bytes" yaml:"-"`
		HashBytes       []byte `json:"-"                    mapstructure:"hash_bytes"       yaml:"-"`
	}

	I18n struct {
		BundleDir string `json:"bundle_dir,omitempty" mapstructure:"bundle_dir" yaml:"bundle_dir,omitempty"`
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

	if cfg.DB.MigrationDir == "" {
		migrationDirPathFromEnv := os.Getenv(constants.MigrationDirPath)
		if migrationDirPathFromEnv != "" {
			cfg.DB.MigrationDir = migrationDirPathFromEnv
		} else {
			cfg.DB.MigrationDir, err = getMigrationDirPath()
			if err != nil {
				return nil, err
			}
		}
	}

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

func getMigrationDirPath() (string, error) {
	// Get the current working directory
	// Getwd gives us the current working directory that we are running our app with `go run` command. in goland we can specify this working directory for the project
	// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Get the absolute path of the executed project directory
	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "migrations/schemas")

	return configPath, nil
}
