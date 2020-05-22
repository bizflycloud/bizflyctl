/*
Copyright © 2020 BizFly Cloud

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
	"log"
	"strings"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	serverListHeader = []string{"ID", "Name", "Status"}
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "BizFly Cloud Server Interaction",
	Long:  `BizFly Cloud Server Action: Create, List, Delete, Resize, Change Type Server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

// deleteCmd represents the delete command
var serverDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Server",
	Long: `Delete Server with server ID as input
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

// listCmd represents the list command
var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all server in your account",
	Long:  `List all server in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := apiClientForContext(cmd)
		servers, err := client.Server.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, server := range servers {
			s := []string{server.ID, server.Name, server.Status}
			data = append(data, s)
		}
		formatter.Output(serverListHeader, data)
	},
}

// getCmd represents the get command
var serverGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a server",
	Long: `Get detail a server with server ID as input
Example: bizfly server get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := apiClientForContext(cmd)

		server, err := client.Server.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Server %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{server.ID, server.Name, server.Status})
		formatter.Output(serverListHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverGetCmd)
	serverCmd.AddCommand(serverDeleteCmd)
}
