package cmd

import (
	"os"

	"github.com/htquangg/a-wasm/internal/cli"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:                   "reset",
	Short:                 "Used to delete all Awasm related data on your machine",
	DisableFlagsInUseLine: true,
	Example:               "awasm reset",
	Args:                  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cli.GetConfigFile()
		if configFile != nil {
			cli.DeleteValueInKeyring(configFile.LoggedInUserEmail)
		}

		_, pathToDir, err := cli.GetFullConfigFilePath()
		if err != nil {
			cli.HandleError(err)
		}

		os.RemoveAll(pathToDir)

		cli.PrintSuccessMessage("Reset successful.")
	},
}

func init() {
}
