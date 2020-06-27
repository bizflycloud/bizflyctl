package cmd

import (
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/spf13/cobra"
)

var imageListHeader = []string{"ID", "Distribution", "Version"}

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "BizFly Cloud Image Interaction",
	Long:  `BizFly Cloud Image Action: List OS Image, Create a custom image`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// imageListcmd represents list all os images
var imageListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all os images in BizFly Cloud",
	Long: `
List all os images in BizFly Cloud
Use: bizfly image list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		images, err := client.Server.ListOSImages(ctx)
		if err != nil {
			fmt.Printf("List os image error: %v", err)
		}
		var data [][]string
		for _, image := range images {
			osDist := image.OSDistribution
			for _, osVer := range image.Version {
				s := []string{osVer.ID, osDist, osVer.Name}
				data = append(data, s)
			}
		}
		formatter.Output(imageListHeader, data)

	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imageListCmd)
}
