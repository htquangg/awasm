package cmd

import "github.com/spf13/cobra"

var endpointsCmd = &cobra.Command{
	Example: "awasm endpoints",
	Use:     "endpoints",
	Short:   "Used to manage endpoints",
	Run:     func(cmd *cobra.Command, args []string) {},
}
