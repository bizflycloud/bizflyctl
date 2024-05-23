package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/spf13/cobra"
)

var (
	flavorListHeader = []string{"ID", "Name", "CPU", "RAM", "Category"}
	vcpus            int
	ram              int
)

// flavorCmd represents the flavor command
var flavorCmd = &cobra.Command{
	Use:   "flavor",
	Short: "Bizfly Cloud Flavor Interaction",
	Long:  `Bizfly Cloud Flavor Action: List Flavors`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// flavorListcmd represents list all flavors
var flavorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flavor of Bizfly Cloud",
	Long: `
List all flavor of Bizfly Cloud.
Use: bizfly flavor list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		flavors, err := client.CloudServer.Flavors().List(ctx)
		if err != nil {
			fmt.Printf("List flavors error %v", err)
			os.Exit(1)
		}
		var data [][]string
		var flavorName string
		for _, flavor := range flavors {
			flavor.RAM = flavor.RAM / 1024
			if category != "" && category != flavor.Category {
				continue
			}
			if vcpus != -1 && vcpus != flavor.VCPUs {
				continue
			}
			if ram != -1 && ram != flavor.RAM {
				continue
			}
			re := regexp.MustCompile(`(\d+c_\d+g)`)
			result := re.FindStringSubmatch(flavor.Name)
			if len(result) == 0 {
				continue
			}
			if result[0] != "" {
				flavorName = result[0]
			}
			s := []string{flavor.ID, flavorName, strconv.Itoa(flavor.VCPUs), strconv.Itoa(flavor.RAM), flavor.Category}
			data = append(data, s)
		}
		formatter.Output(flavorListHeader, data)

	},
}

func init() {
	rootCmd.AddCommand(flavorCmd)
	flavorCmd.AddCommand(flavorListCmd)
	flpf := flavorListCmd.PersistentFlags()
	flpf.StringVar(&category, "category", "", "Filter flavor by category")
	flpf.IntVar(&vcpus, "cpu", -1, "Filter flavor by cpus")
	flpf.IntVar(&ram, "ram", -1, "Filter flavor by ram")
}
