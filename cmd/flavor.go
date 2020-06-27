package cmd

import (
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/spf13/cobra"
	"os"
)

var flavorListHeader = []string{"ID", "Name"}

// flavorCmd represents the flavor command
var flavorCmd = &cobra.Command{
	Use:   "flavor",
	Short: "BizFly Cloud Flavor Interaction",
	Long:  `BizFly Cloud Flavor Action: List Flavors`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// flavorListcmd represents list all flavors
var flavorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flavor of BizFly Cloud",
	Long: `
List all flavor of BizFly Cloud.
Use: bizfly flavor list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		flavors, err := client.Server.ListFlavors(ctx)
		if err != nil {
			fmt.Printf("List flavors error %v", err)
			os.Exit(1)
		}
		var data [][]string
		for _, flavor := range flavors {
			s := []string{flavor.ID, flavor.Name}
			data = append(data, s)
		}
		formatter.Output(flavorListHeader, data)

	},
}

func init() {
	rootCmd.AddCommand(flavorCmd)
	flavorCmd.AddCommand(flavorListCmd)
}
