package constants

import "time"

const (
	PrefixEnv        = "AWASM"
	ConfigPathEnv    = PrefixEnv + "_" + "CONFIG_PATH"
	MigrationDirPath = PrefixEnv + "_" + "DB_MIGRATION_DIR_PATH"
	LogLevelEnv      = PrefixEnv + "_" + "LOG_LEVEL"
	LogPathEnv       = PrefixEnv + "_" + "LOG_PATH"

	Yaml = "yaml"

	Dev        = "development"
	Test       = "test"
	Production = "production"

	ReadTimeout       = 10 * time.Minute
	WriteTimeout      = 60 * time.Second
	ReadHeaderTimeout = 30 * time.Second
	ShutdownTimeout   = 5 * time.Second
)

var (
	Version   = ""
	Revision  = ""
	GoVersion = ""
)
