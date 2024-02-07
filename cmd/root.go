package cmd

import (
	"os"
	"strings"

	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "awasm",
	Short: "Awasm is the application that you can build, deploy, and run your application on the edge.",
	Long: `Awasm is the application that you can build, deploy, and run your application on the edge.
To run awasm, use:
  - 'awasm run' to launch application.
  - 'awasm endpoints' to manage endpoints.
		`,
}

func Execute() {
	initLog()

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("the service exitted abnormally")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog)

	for _, cmd := range []*cobra.Command{runCmd, endpointsCmd} {
		rootCmd.AddCommand(cmd)
	}
}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logLevel := os.Getenv(constants.LogLevel)
	switch strings.ToLower(logLevel) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
