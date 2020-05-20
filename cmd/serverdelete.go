/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/bizflycloud/gobizfly"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var serverDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Server",
	Long: `Get detail a server with server ID as input
Example: bizfly server delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58

You can delete multiple server with list of server id
Example: bizfly server delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58 f5869e9c-9ab2-11ea-b9e3-e353a4f04836
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := apiClientForContext(cmd)
		for _, serverID := range args {
			fmt.Printf("Deleting server %s \n", serverID)
			err := client.Server.Delete(ctx, serverID)
			if err != nil {
				if errors.Is(err, gobizfly.ErrNotFound) {
					fmt.Printf("Server %s is not found", serverID)
					return
				}
			}
		}
	},
}

func init() {
	serverCmd.AddCommand(serverDeleteCmd)
}
