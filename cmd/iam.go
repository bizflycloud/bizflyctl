package cmd

import (
	"fmt"
	"strconv"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/spf13/cobra"
)

var projectListHeader = []string{"ID", "Name", "Description", "Is Active", "Created At", "Updated At"}

var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Bizfly Cloud IAM Interaction",
	Long:  `Bizfly Cloud IAM Action: List Projects`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iam called")
	},
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Bizfly Cloud Projects Interaction",
	Long:  `Bizfly Cloud Projects Action: List Projects`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("projects called")
	},
}
var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in Bizfly Cloud",
	Long: `
List all projects in Bizfly Cloud
Use: bizfly projects list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		projects, err := client.IAM.ListProjects(ctx)
		if err != nil {
			fmt.Printf("List projects error: %v", err)
		}
		var data [][]string
		for _, project := range projects {
			s := []string{project.UUID, project.AliasName, project.Description, strconv.FormatBool(project.IsActive),
				project.CreatedAt, project.UpdatedAt}
			data = append(data, s)
		}
		formatter.Output(projectListHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(projectsListCmd)
}
