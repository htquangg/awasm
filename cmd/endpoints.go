package cmd

import (
	"fmt"
	"os"

	"github.com/htquangg/a-wasm/internal/cli"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
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
	publishEndpointCmd.Flags().String("deployment-id", "", "The id of the deployment that you want to publish LIVE")
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

		loggedInUserDetails, err := cli.GetCurrentLoggedInUserDetails()
		if err != nil {
			cli.HandleError(err, "Unable to authenticate")
		}

		// set up resty client
		httpClient := resty.New()
		httpClient.
			SetAuthToken(loggedInUserDetails.UserCredentials.AccessToken).
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json")

		addEndpointResp, err := api.CallAddEndpoint(httpClient, &schemas.AddEndpointReq{
			Name:    name,
			Runtime: runtime,
		})
		if err != nil {
			cli.HandleError(err)
		}

		fmt.Println("Successful!!!")
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

		loggedInUserDetails, err := cli.GetCurrentLoggedInUserDetails()
		if err != nil {
			cli.HandleError(err, "Unable to authenticate")
		}

		// set up resty client
		httpClient := resty.New()
		httpClient.
			SetAuthToken(loggedInUserDetails.UserCredentials.AccessToken).
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json")

		publicEndpointResp, err := api.CallPublishEndpoint(httpClient, &schemas.PublishEndpointReq{
			DeploymentID: deploymentID,
		})
		if err != nil {
			cli.HandleError(err)
		}

		fmt.Println("Successful!!!")
		fmt.Printf("Your ingress url: %s\n", publicEndpointResp.IngressURL)
	},
}
