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
	"github.com/bizflycloud/bizflycli/common"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

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
		common.Output(serverListHeader, data)
	},
}

func init() {
	serverCmd.AddCommand(serverGetCmd)
}
