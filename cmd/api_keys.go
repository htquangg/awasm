package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/htquangg/a-wasm/internal/cli"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"
)

var apiKeysCmd = &cobra.Command{
	Example:               "awasm apikeys",
	Use:                   "apikeys",
	Short:                 "Used to manage api keys",
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.CheckAuthentication()
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	createApiKeyCmd.Flags().String("name", "", "The friendly name of api key")
	_ = createApiKeyCmd.MarkFlagRequired("name")
	apiKeysCmd.AddCommand(createApiKeyCmd)
}

var createApiKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Used to create new api key",
	Run: func(cmd *cobra.Command, args []string) {
		friendlyName, err := cmd.Flags().GetString("name")
		if err != nil {
			cli.PrintErrorMessageAndExit("unable to parse flag name. --name <friendly-name>")
		}

		loggedInUserDetails, isAuthenticated, err := cli.GetCurrentLoggedInUserDetails()
		if err != nil {
			cli.HandleError(err, "Unable to authenticate")
		}
		if !isAuthenticated {
			cli.PrintErrorMessageAndExit(
				"Your login session has expired, please run [awasm login] and try again",
			)
		}

		client := api.NewClient(&api.ClientOptions{
			Debug: viper.GetBool("cli.debug"),
		})
		client.HTTPClient.SetAuthToken(loggedInUserDetails.UserCredentials.AccessToken)

		addApiKeyResp, err := api.CallAddApiKey(client.HTTPClient, &schemas.AddApiKeyReq{
			FriendlyName: friendlyName,
		},
		)
		if err != nil {
			cli.HandleError(err)
		}

		fmt.Printf("Your api key: %s\n", addApiKeyResp.Key)
	},
}
