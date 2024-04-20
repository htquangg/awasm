package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/pkg/logger"

	"github.com/spf13/cobra"
)

// go build -ldflags "-X github.com/htquangg/a-wasm/cmd.Version=x.y.z"
var (
	Name      = "awasm"
	Version   = "devel"
	Revision  = ""
	Time      = ""
	GoVersion = "1.21"
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
	Version:           fmt.Sprintf("%s\nrevision: %s\nbuild time: %s", Version, Revision, Time),
}

// exitCode wraps a return value for the application
type exitCode struct {
	Err  error
	Code int
}

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(exitCode); ok {
			if exit.Code != 0 {
				fmt.Fprintln(
					os.Stderr,
					"Awasm failed at",
					time.Now().Format("January 2, 2006 at 3:04pm (MST)"),
					"Err:",
					exit.Err,
				)
			} else {
				fmt.Fprintln(os.Stderr, "Stopped Awasm at", time.Now().Format("January 2, 2006 at 3:04pm (MST)"))
			}

			os.Exit(exit.Code)
		}
		panic(e) // not an exitCode, bubble up
	}
}

func Execute() {
	// This makes sure that we panic and run defers correctly
	defer handleExit()

	initLog()

	rootCmd.PersistentFlags().
		StringVar(&config.AWASM_URL, "domain", constants.AWASM_DEFAULT_API_URL, "Point the CLI to your own backend [can also set via environment variable name: AWASM_API_URL]")

	// if config.AWASM_URL is set to the default value, check if AWASM_URL is set in the environment
	// this is used to allow overrides of the default value
	if !rootCmd.Flag("domain").Changed {
		if envAwasmBackendUrl, ok := os.LookupEnv("AWASM_API_URL"); ok {
			config.AWASM_URL = envAwasmBackendUrl
		}
	}

	if err := rootCmd.Execute(); err != nil {
		panic(exitCode{
			Code: 1,
			Err:  err,
		})
	}

	panic(exitCode{Code: 0})
}

func init() {
	cobra.OnInitialize(initLog)

	for _, cmd := range []*cobra.Command{runCmd, endpointsCmd, deploymentsCmd, loginCmd, signupCmd, resetCmd} {
		rootCmd.AddCommand(cmd)
	}
}

func initLog() {
	logPath := os.Getenv(constants.LogPath)
	logLevel := os.Getenv(constants.LogLevel)

	logger.SetLogger(logger.NewZapLogger(
		logger.WithZapFilename(logPath),
		logger.WithZapLevel(logLevel),
	))
}
