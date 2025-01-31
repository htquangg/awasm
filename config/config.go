package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/htquangg/awasm/internal/constants"
	"github.com/htquangg/awasm/pkg/converter"
	"github.com/htquangg/awasm/pkg/logger"
)

var AwasmUrl string

type MailerProviderType string

const (
	ProviderTypeNoop MailerProviderType = "NOOP"
	ProviderTypeSMTP MailerProviderType = "SIMPLE_SMTP"
)

type (
	Config struct {
		HTTP       *HTTP     `json:"http,omitempty"       mapstructure:"http"        yaml:"http,omitempty"`
		DB         *DB       `json:"db,omitempty"         mapstructure:"db"          yaml:"db,omitempty"`
		Redis      *Redis    `json:"redis,omitempty"      mapstructure:"redis"       yaml:"redis,omitempty"`
		JWT        *JWT      `json:"jwt,omitempty"        mapstructure:"jwt"         yaml:"jwt,omitempty"`
		Key        *Key      `json:"key,omitempty"        mapstructure:"key"         yaml:"key,omitempty"`
		Mailer     *Mailer   `json:"mailer"               mapstructure:"mailer"      yaml:"mailer"`
		Logging    *Logging  `json:"logging"              mapstructure:"logging"     yaml:"logging"`
		Session    *Session  `json:"session"              mapstructure:"session"     yaml:"session"`
		Security   *Security `json:"security"             mapstructure:"security"    yaml:"security"`
		I18n       *I18n     `json:"i18n,omitempty"       mapstructure:"i18n"        yaml:"i18n,omitempty"`
		IngressURL string    `json:"ingressURL,omitempty" mapstructure:"ingress_url" yaml:"ingress_url,omitempty"`
	}

	HTTP struct {
		Addr            string `json:"addr"            mapstructure:"addr"              yaml:"addr"`
		ShowStartBanner bool   `json:"showStartBanner" mapstructure:"show_start_banner" yaml:"show_start_banner"`
	}

	DB struct {
		Host             string `json:"host,omitempty"     mapstructure:"host"               yaml:"host,omitempty"`
		User             string `json:"user,omitempty"     mapstructure:"user"               yaml:"user,omitempty"`
		Password         string `json:"password,omitempty" mapstructure:"password"           yaml:"password,omitempty"`
		Schema           string `json:"schema,omitempty"   mapstructure:"schema"             yaml:"schema,omitempty"`
		Charset          string `json:"charset"            mapstructure:"charset"            yaml:"charset"`
		SslMode          string `json:"sslMode"            mapstructure:"ssl_mode"           yaml:"ssl_mode"`
		MigrationDirPath string `json:"migrationDirPath"   mapstructure:"migration_dir_path" yaml:"migration_dir_path"`
		MaxIdleConns     int    `json:"maxIdleConns"       mapstructure:"max_idle_conns"     yaml:"max_idle_conns"`
		MaxOpenConns     int    `json:"maxOpenConns"       mapstructure:"max_open_conns"     yaml:"max_open_conns"`
		ConnMaxLifetime  int    `json:"connMaxLifetime"    mapstructure:"conn_max_lifetime"  yaml:"conn_max_lifetime"`
		Port             int    `json:"port,omitempty"     mapstructure:"port"               yaml:"port,omitempty"`
		LogSQL           bool   `json:"logSql"             mapstructure:"log_sql"            yaml:"log_sql"`
	}

	Redis struct {
		Host       string `json:"host,omitempty" mapstructure:"host"        yaml:"host,omitempty"`
		Password   string `json:"password"       mapstructure:"password"    yaml:"password"`
		Port       int    `json:"port,omitempty" mapstructure:"port"        yaml:"port,omitempty"`
		DB         int    `json:"db,omitempty"   mapstructure:"db"          yaml:"db,omitempty"`
		PoolSize   int    `json:"poolSize"       mapstructure:"pool_size"   yaml:"pool_size"`
		RequireTLS bool   `json:"requireTLS"     mapstructure:"require_tls" yaml:"require_tls"`
	}

	JWT struct {
		Secret      string `json:"secret,omitempty" mapstructure:"secret"       yaml:"secret,omitempty"`
		SecretBytes []byte `json:"-"                mapstructure:"secret_bytes" yaml:"-"`
		Exp         int    `json:"exp,omitempty"    mapstructure:"exp"          yaml:"exp,omitempty"`
	}

	Key struct {
		Encryption               string `json:"encryption,omitempty"          mapstructure:"encryption"                   yaml:"encryption,omitempty"`
		Hash                     string `json:"hash,omitempty"                mapstructure:"hash"                         yaml:"hash,omitempty"`
		ApiKeySignatureHMAC      string `json:"apiKeySignatureHmac,omitempty" mapstructure:"api_key_signature_hmac"       yaml:"api_key_signature_hmac,omitempty"`
		ApiKeyDatabaseHMAC       string `json:"apiKeyDatabaseHmac,omitempty"  mapstructure:"api_key_database_hmac"        yaml:"api_key_database_hmac,omitempty"`
		CacheKeyHMAC             string `json:"cacheKeyHmac,omitempty"        mapstructure:"cache_key_hmac"               yaml:"cache_key_hmac,omitempty"`
		EncryptionBytes          []byte `json:"-"                             mapstructure:"encryption_bytes"             yaml:"-"`
		HashBytes                []byte `json:"-"                             mapstructure:"hash_bytes"                   yaml:"-"`
		ApiKeySignatureHMACBytes []byte `json:"-"                             mapstructure:"api_key_signature_hmac_bytes" yaml:"-"`
		ApiKeyDatabaseHMACBytes  []byte `json:"-"                             mapstructure:"api_key_database_hmac_bytes"  yaml:"-"`
		CacheKeyHMACBytes        []byte `json:""                              mapstructure:"cache_key_hmac_bytes"         yaml:""`
	}

	Mailer struct {
		ProviderType MailerProviderType `json:"providerType,omitempty" mapstructure:"provider_type" yaml:"provider_type,omitempty"`
		FromEmail    string             `json:"fromEmail,omitempty"    mapstructure:"from_email"    yaml:"from_email,omitempty"`
		FromName     string             `json:"fromName"               mapstructure:"from_name"     yaml:"from_name"`
		User         string             `json:"user"                   mapstructure:"user"          yaml:"user"`
		Password     string             `json:"password"               mapstructure:"password"      yaml:"password"`
		Host         string             `json:"host,omitempty"         mapstructure:"host"          yaml:"host,omitempty"`
		Port         int                `json:"port,omitempty"         mapstructure:"port"          yaml:"port,omitempty"`
		RequireTLS   bool               `json:"requireTLS"             mapstructure:"require_tls"   yaml:"require_tls"`
	}

	Logging struct {
		Filename     string `json:"filename"     mapstructure:"filename"       yaml:"filename"`
		Level        string `json:"level"        mapstructure:"level"          yaml:"level"`
		MaxBackups   int    `json:"maxBackups"   mapstructure:"max_backups"    yaml:"max_backups"`
		MaxSize      int    `json:"maxSize"      mapstructure:"max_size"       yaml:"max_size"`
		MaxAge       int    `json:"maxAge"       mapstructure:"max_age"        yaml:"max_age"`
		UseLocalTime bool   `json:"useLocalTime" mapstructure:"use_local_time" yaml:"use_local_time"`
		UseCompress  bool   `json:"useCompress"  mapstructure:"use_compress"   yaml:"use_compress"`
	}

	Session struct {
		Timebox           *time.Duration `json:"timebox"                     mapstructure:"timebox"            yaml:"timebox"`
		InactivityTimeout *time.Duration `json:"inactivityTimeout,omitempty" mapstructure:"inactivity_timeout" yaml:"inactivity_timeout,omitempty"`
	}

	Security struct {
		RefreshTokenReuseInterval int `json:"refreshTokenReuseInterval" mapstructure:"refresh_token_reuse_interval" yaml:"refresh_token_reuse_interval"`
	}

	I18n struct {
		BundleDir string `json:"bundleDir,omitempty" mapstructure:"bundle_dir" yaml:"bundle_dir,omitempty"`
	}
)

func LoadConfig() (*Config, error) {
	loadDefaultConfig()

	// Load the configuration from the file
	var configPath string

	// TOIMPROVE:
	// Will remove configuration from env, if i can find the way to inject configuration to the unit test.
	configPathFromEnv := os.Getenv(constants.ConfigPathEnv)
	if configPathFromEnv != "" {
		configPath = configPathFromEnv
	} else if viper.GetString("server.config-path") != "" {
		configPath = viper.GetString("server.config-path")
	} else {
		rootPath, err := getConfigRootPath()
		if err != nil {
			return nil, err
		}
		configPath = fmt.Sprintf("%s/awasm.yaml", rootPath)
	}

	logger.Infof("Load config from %s\n", configPath)

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// setup viper to be able to read env variables with a configured prefix
	viper.SetEnvPrefix(constants.PrefixEnv)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

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

	apiKeySignatureHMACBytes, err := converter.FromURLB64(cfg.Key.ApiKeySignatureHMAC)
	if err != nil {
		return nil, fmt.Errorf("could not decode api-key-signature-hmac: %v", err)
	}
	cfg.Key.ApiKeySignatureHMACBytes = apiKeySignatureHMACBytes

	apiKeyDatabaseHMACBytes, err := converter.FromURLB64(cfg.Key.ApiKeyDatabaseHMAC)
	if err != nil {
		return nil, fmt.Errorf("could not decode api-key-database-hmac: %v", err)
	}
	cfg.Key.ApiKeyDatabaseHMACBytes = apiKeyDatabaseHMACBytes

	cacheKeyHMACBytes, err := converter.FromB64(cfg.Key.CacheKeyHMAC)
	if err != nil {
		return nil, fmt.Errorf("could not decode cache-key-hmac: %v", err)
	}
	cfg.Key.CacheKeyHMACBytes = cacheKeyHMACBytes

	jwtSecretBytes, err := converter.FromURLB64(cfg.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("could not decode jwt-secret-key: %v", err)
	}
	cfg.JWT.SecretBytes = jwtSecretBytes

	if cfg.DB.MigrationDirPath == "" {
		migrationDirPathFromEnv := os.Getenv(constants.MigrationDirPath)
		if migrationDirPathFromEnv != "" {
			cfg.DB.MigrationDirPath = migrationDirPathFromEnv
		} else {
			cfg.DB.MigrationDirPath, err = getMigrationDirPath()
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

	viper.SetDefault("db.charset", "utf8bm4")
	viper.SetDefault("db.log_sql", true)
	viper.SetDefault("db.ssl_mode", "disable")
	viper.SetDefault("db.max_idle_conns", 100)
	viper.SetDefault("db.max_open_conns", 100)
	viper.SetDefault("db.conn_max_lifetime", 300)

	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.pool_size", 50)
	viper.SetDefault("redis.require_tls", false)

	viper.SetDefault("mailer.provider_type", ProviderTypeNoop)
	viper.SetDefault("mailer.from_name", "Local Awasm")
	viper.SetDefault("mailer.user", "")
	viper.SetDefault("mailer.password", "")
	viper.SetDefault("mailer.require_tls", false)

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
