package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/htquangg/a-wasm/internal/cli"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"
)

var endpointsCmd = &cobra.Command{
	Example:               "awasm endpoints",
	Use:                   "endpoints",
	Short:                 "Used to manage endpoints",
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cli.CheckAuthentication()
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	// add createEndpointCmd flag here
	createEndpointCmd.Flags().String("name", "", "The name of your endpoint")
	_ = createEndpointCmd.MarkFlagRequired("name")
	createEndpointCmd.Flags().String("runtime", "", "The runtime of your endpoint (go or js)")
	_ = createEndpointCmd.MarkFlagRequired("runtime")

	// add publishEndpointCmd flag here
	publishEndpointCmd.Flags().
		String("deployment-id", "", "The id of the deployment that you want to publish LIVE")
	_ = publishEndpointCmd.MarkFlagRequired("deployment-id")

	endpointsCmd.AddCommand(createEndpointCmd, publishEndpointCmd)
}

var createEndpointCmd = &cobra.Command{
	Use:   "create",
	Short: "Used to create new endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("unable to parse flag name. --name <name>")
			os.Exit(1)
		}

		runtime, err := cmd.Flags().GetString("runtime")
		if err != nil {
			fmt.Println("unable to parse flag runtime [--runtime go, --runtime js]")
			os.Exit(1)
		}
		if !schemas.ValidRuntime(runtime) {
			fmt.Printf("invalid runtime %s, only go and js are currently supported\n", runtime)
			os.Exit(1)
		}

		loggedInUserDetails, isAuthenticated, err := cli.GetCurrentLoggedInUserDetails()
		if err != nil {
			cli.HandleError(err, "unable to authenticate")
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

		addEndpointResp, err := api.CallAddEndpoint(client.HTTPClient, &schemas.AddEndpointReq{
			Name:    name,
			Runtime: runtime,
		})
		if err != nil {
			cli.HandleError(err)
		}

		fmt.Printf("Your endpoint id: %s\n", addEndpointResp.ID)
	},
}

var publishEndpointCmd = &cobra.Command{
	Use:   "publish",
	Short: "Used to publish deployment to an endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		deploymentID, err := cmd.Flags().GetString("deployment-id")
		if err != nil {
			fmt.Println("unable to parse flag deployment id. --deployment-id <deploymentID>")
			os.Exit(1)
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

		publicEndpointResp, err := api.CallPublishEndpoint(
			client.HTTPClient,
			&schemas.PublishEndpointReq{
				DeploymentID: deploymentID,
			},
		)
		if err != nil {
			cli.HandleError(err)
		}

		fmt.Printf("Your ingress url: %s\n", publicEndpointResp.IngressURL)
	},
}
