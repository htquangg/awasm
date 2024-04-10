package cmd

import (
	"fmt"
	"os"

	"github.com/htquangg/a-wasm/internal/cli"
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
	PreRun: func(cmd *cobra.Command, args []string) {
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

		// set up resty client
		httpClient := resty.New()
		httpClient.
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json")

		body := schemas.AddEndpointReq{
			Name:    name,
			Runtime: runtime,
		}

		httpRequest := httpClient.
			R().
			SetBody(body)

		response, err := httpRequest.Post(fmt.Sprintf("%v/api/v1/endpoints/", "http://127.0.0.1:3000"))
		if err != nil {
			fmt.Printf("CallCreateEndpointV1: Unable to complete api request [err=%s]\n", err)
			os.Exit(1)
		}
		if response.IsError() {
			fmt.Printf("CallCreateEndpointV1: Unsuccessful [response=%s]\n", response.String())
			os.Exit(1)
		}

		fmt.Printf("CallCreateEndpointV1: Successful [response=%s]\n", response.String())
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

		// set up resty client
		httpClient := resty.New()
		httpClient.
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json")

		body := schemas.PublishEndpointReq{
			DeploymentID: deploymentID,
		}

		httpRequest := httpClient.
			R().
			SetBody(body)

		response, err := httpRequest.Post(fmt.Sprintf("%v/api/v1/live/publish", "http://127.0.0.1:3000"))
		if err != nil {
			fmt.Printf("CallPublishEndpointV1: Unable to complete api request [err=%s]\n", err)
			os.Exit(1)
		}
		if response.IsError() {
			fmt.Printf("CallPublishEndpointV1: Unsuccessful [response=%s]\n", response.String())
			os.Exit(1)
		}

		fmt.Printf("CallPublishEndpointV1: Successful [response=%s]\n", response.String())
	},
}
