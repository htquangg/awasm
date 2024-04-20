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

var deploymentsCmd = &cobra.Command{
	Example:               "awasm deployments",
	Use:                   "deployments",
	Short:                 "Used to manage deployments",
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		cli.CheckAuthentication()
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	// add createDeploymentCmd flag here
	createDeploymentCmd.Flags().String("endpoint-id", "", "The id of the endpoint to where you want to deploy")
	_ = createDeploymentCmd.MarkFlagRequired("endpoint-id")
	createDeploymentCmd.Flags().String("file", "", "The file location of your code that you want to deploy")
	_ = createDeploymentCmd.MarkFlagRequired("file")
	deploymentsCmd.AddCommand(createDeploymentCmd)
}

var createDeploymentCmd = &cobra.Command{
	Use:   "create",
	Short: "Used to create new deployment",
	Run: func(cmd *cobra.Command, args []string) {
		endpointID, err := cmd.Flags().GetString("endpoint-id")
		if err != nil {
			fmt.Println("unable to parse flag endpoint id. --endpoint-id <deploymentID>")
			os.Exit(1)
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("unable to parse flag file. --file <filename>")
		}

		b, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("unable to read file")
		}

		loggedInUserDetails, isAuthenticated, err := cli.GetCurrentLoggedInUserDetails()
		if err != nil {
			cli.HandleError(err, "Unable to authenticate")
		}
		if !isAuthenticated {
			cli.PrintErrorMessageAndExit("Your login session has expired, please run [awasm login] and try again")
		}

		// set up resty client
		httpClient := resty.New()
		httpClient.
			SetAuthToken(loggedInUserDetails.UserCredentials.AccessToken).
			SetHeader("Accept", "application/octet-stream").
			SetHeader("Content-Type", "application/octet-stream")

		addDeploymentResp, err := api.CallAddDeployment(httpClient, &schemas.AddDeploymentReq{
			EndpointID: endpointID,
			Data:       b,
		},
		)
		if err != nil {
			cli.HandleError(err)
		}

		fmt.Println("Successful!!!")
		fmt.Printf("Your deployment id: %s\n", addDeploymentResp.ID)
		fmt.Printf("Your preview ingress url: %s\n", addDeploymentResp.IngressURL)
	},
}
