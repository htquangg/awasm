package cmd

import (
	"os"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/segmentfault/pacman/contrib/log/zap"
	"github.com/segmentfault/pacman/log"
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
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	Version:           constants.CLI_VERSION,
}

func Execute() {
	initLog()

	rootCmd.PersistentFlags().
		StringVar(&config.AWASM_URL, "domain", constants.AWASM_DEFAULT_API_URL, "Point the CLI to your own backend [can also set via environment variable name: AWASM_API_URL]")

	// if config.AWASM_URL is set to the default value, check if AWASM_URL is set in the environment
	// this is used to allow overrides of the default value
	if !rootCmd.Flag("domain").Changed {
		if envInfisicalBackendUrl, ok := os.LookupEnv("AWASM_API_URL"); ok {
			config.AWASM_URL = envInfisicalBackendUrl
		}
	}

	if err := rootCmd.Execute(); err != nil {
		log.Errorf("the service exitted abnormally: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog)

	for _, cmd := range []*cobra.Command{runCmd, endpointsCmd, deploymentsCmd, loginCmd, signupCmd, resetCmd} {
		rootCmd.AddCommand(cmd)
	}
}

func initLog() {
	logLevel := os.Getenv(constants.LogLevel)
	logPath := os.Getenv(constants.LogPath)

	log.SetLogger(zap.NewLogger(
		log.ParseLevel(logLevel), zap.WithName("a-wasm"), zap.WithPath(logPath)))
}
