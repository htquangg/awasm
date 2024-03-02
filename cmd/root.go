package cmd

import (
	"os"

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
}

func Execute() {
	initLog()

	if err := rootCmd.Execute(); err != nil {
		log.Errorf("the service exitted abnormally: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog)

	for _, cmd := range []*cobra.Command{runCmd, endpointsCmd, deploymentsCmd} {
		rootCmd.AddCommand(cmd)
	}
}

func initLog() {
	logLevel := os.Getenv(constants.LogLevel)
	logPath := os.Getenv(constants.LogPath)

	log.SetLogger(zap.NewLogger(
		log.ParseLevel(logLevel), zap.WithName("a-wasm"), zap.WithPath(logPath)))
}
