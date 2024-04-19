package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/pkg/converter"

	"github.com/spf13/viper"
)

var AWASM_URL string

type (
	Config struct {
		Server     *Server  `json:"server,omitempty"     mapstructure:"server"      yaml:"server,omitempty"`
		DB         *DB      `json:"db,omitempty"         mapstructure:"db"          yaml:"db,omitempty"`
		Redis      *Redis   `json:"redis,omitempty"      mapstructure:"redis"       yaml:"redis,omitempty"`
		JWT        *JWT     `json:"jwt,omitempty"        mapstructure:"jwt"         yaml:"jwt,omitempty"`
		Key        *Key     `json:"key,omitempty"        mapstructure:"key"         yaml:"key,omitempty"`
		I18n       *I18n    `json:"i18n,omitempty"       mapstructure:"i18n"        yaml:"i18n,omitempty"`
		Logging    *Logging `json:"logging"              mapstructure:"logging"     yaml:"logging"`
		IngressURL string   `json:"ingressURL,omitempty" mapstructure:"ingress_url" yaml:"ingress_url,omitempty"`
	}

	Server struct {
		Addr            string `json:"addr,omitempty"  mapstructure:"addr"              yaml:"addr,omitempty"`
		ShowStartBanner bool   `json:"showStartBanner" mapstructure:"show_start_banner" yaml:"show_start_banner"`
	}

	DB struct {
		Host            string `json:"host,omitempty"     mapstructure:"host"              yaml:"host,omitempty"`
		User            string `json:"user,omitempty"     mapstructure:"user"              yaml:"user,omitempty"`
		Password        string `json:"password,omitempty" mapstructure:"password"          yaml:"password,omitempty"`
		Schema          string `json:"schema,omitempty"   mapstructure:"schema"            yaml:"schema,omitempty"`
		Charset         string `json:"charset"            mapstructure:"charset"           yaml:"charset"`
		SslMode         string `json:"sslMode"            mapstructure:"ssl_mode"          yaml:"ssl_mode"`
		MigrationDir    string `json:"migrationDir"       mapstructure:"migration_dir"     yaml:"migration_dir"`
		MaxIdleConns    int    `json:"maxIdleConns"       mapstructure:"max_idle_conns"    yaml:"max_idle_conns"`
		MaxOpenConns    int    `json:"maxOpenConns"       mapstructure:"max_open_conns"    yaml:"max_open_conns"`
		ConnMaxLifetime int    `json:"connMaxLifetime"    mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
		Port            uint16 `json:"port,omitempty"     mapstructure:"port"              yaml:"port,omitempty"`
		LogSQL          bool   `json:"logSql"             mapstructure:"log_sql"           yaml:"log_sql"`
	}

	Redis struct {
		Host       string `json:"host,omitempty"       mapstructure:"host"        yaml:"host,omitempty"`
		Password   string `json:"password"             mapstructure:"password"    yaml:"password"`
		Port       int    `json:"port,omitempty"       mapstructure:"port"        yaml:"port,omitempty"`
		DB         int    `json:"db,omitempty"         mapstructure:"db"          yaml:"db,omitempty"`
		PoolSize   int    `json:"poolSize"             mapstructure:"pool_size"   yaml:"pool_size"`
		RequireTLS bool   `json:"requireTLS,omitempty" mapstructure:"require_tls" yaml:"require_tls,omitempty"`
	}

	JWT struct {
		Secret      string `json:"secret,omitempty" mapstructure:"secret"       yaml:"secret,omitempty"`
		SecretBytes []byte `json:"-"                mapstructure:"secret_bytes" yaml:"-"`
		Exp         int    `json:"exp,omitempty"    mapstructure:"exp"          yaml:"exp,omitempty"`
	}

	Key struct {
		Encryption      string `json:"encryption,omitempty" mapstructure:"encryption"       yaml:"encryption,omitempty"`
		Hash            string `json:"hash,omitempty"       mapstructure:"hash"             yaml:"hash,omitempty"`
		EncryptionBytes []byte `json:"-"                    mapstructure:"encryption_bytes" yaml:"-"`
		HashBytes       []byte `json:"-"                    mapstructure:"hash_bytes"       yaml:"-"`
	}

	Logging struct {
		Filename     string `json:"filename"   mapstructure:"filename"    yaml:"filename"`
		Level        string `json:"level"      mapstructure:"level"       yaml:"level"`
		MaxBackups   int    `json:"maxBackups" mapstructure:"max_backups" yaml:"max_backups"`
		MaxSize      int    `json:"maxSize"    mapstructure:"max_size"    yaml:"max_size"`
		MaxAge       int    `json:"maxAge"     mapstructure:"max_age"     yaml:"max_age"`
		UseLocalTime bool   `json:"useLocalTime"  mapstructure:"use_local_time"  yaml:"use_local_time"`
		UseCompress  bool   `json:"useCompress"   mapstructure:"use_compress"    yaml:"use_compress"`
	}

	I18n struct {
		BundleDir string `json:"bundle_dir,omitempty" mapstructure:"bundle_dir" yaml:"bundle_dir,omitempty"`
	}
)

func LoadConfig() (*Config, error) {
	loadDefaultConfig()

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

	secretEncryptionKeyBytes, err := converter.FromB64(cfg.Key.Encryption)
	if err != nil {
		return nil, fmt.Errorf("could not decode email-encryption-key: %v", err)
	}
	cfg.Key.EncryptionBytes = secretEncryptionKeyBytes

	hashingKeyBytes, err := converter.FromB64(cfg.Key.Hash)
	if err != nil {
		return nil, fmt.Errorf("could not decode email-hash-key: %v", err)
	}
	cfg.Key.HashBytes = hashingKeyBytes

	jwtSecretBytes, err := converter.FromURLB64(cfg.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("could not decode email-hash-key: %v", err)
	}
	cfg.JWT.SecretBytes = jwtSecretBytes

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

func loadDefaultConfig() {
	viper.SetDefault("ingress_url", "http://127.0.0.1:8080")

	viper.SetDefault("server.addr", "127.0.0.1:8080")
	viper.SetDefault("server.show_start_banner", true)

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "localdb")
	viper.SetDefault("db.schema", "dev-local-awasm-001")
	viper.SetDefault("db.charset", "utf8bm4")
	viper.SetDefault("db.log_sql", true)
	viper.SetDefault("db.ssl_mode", "disable")
	viper.SetDefault("db.max_idle_conns", 100)
	viper.SetDefault("db.max_open_conns", 100)
	viper.SetDefault("db.conn_max_lifetime", 300)

	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.require_tls", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 1)
	viper.SetDefault("redis.pool_size", 50)

	viper.SetDefault("jwt.secret", "lLTQ-5romlnVEOVdq6gCwEpCCKeypoKvuugm2GDfjzs=")
	viper.SetDefault("jwt.exp", 86400)

	viper.SetDefault("key.encryption", "bASFEEq6OmFvrpYDGgsF7lZn95p8VkuDgPAw93hpad8=")
	viper.SetDefault("key.hash", "lXdxDsEADpazx2V9vR6tjnDa60CVdEaIp2z8jLChTR0oyqpcWm0fQcDq7dKAzq44ttGcN90TvjmsC67llMsz8w==")

	viper.SetDefault("logging.filename", "logs/awasm.log")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.max_backups", 30)
	viper.SetDefault("logging.max_age", 10)
	viper.SetDefault("logging.use_local_time", false)
	viper.SetDefault("logging.use_compress", true)

	viper.SetDefault("i18n.bundle_dir", "./i18n")
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
