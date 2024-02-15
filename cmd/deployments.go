package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var deploymentsCmd = &cobra.Command{
	Example: "awasm deployments",
	Use:     "deployments",
	Short:   "Used to manage deployments",
	Run:     func(cmd *cobra.Command, args []string) {},
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
			fmt.Println("unable to parse flag endpoint id. --endpoint-id <uuid>")
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

		// set up resty client
		httpClient := resty.New()
		httpClient.
			SetHeader("Accept", "application/octet-stream").
			SetHeader("Content-Type", "application/octet-stream")

		httpRequest := httpClient.
			R().
			SetBody(b)

		response, err := httpRequest.Post(fmt.Sprintf("%v/api/v1/endpoints/%s/deployments", "http://127.0.0.1:3000", endpointID))
		if err != nil {
			fmt.Printf("CallCreateDeploymentV1: Unable to complete api request [err=%s]\n", err)
			os.Exit(1)
		}
		if response.IsError() {
			fmt.Printf("CallCreateDeploymentV1: Unsuccessful [response=%s]\n", response.String())
			os.Exit(1)
		}

		fmt.Printf("CallCreateDeploymentV1: Successful [response=%s]\n", response.String())
	},
}
