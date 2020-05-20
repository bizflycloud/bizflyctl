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
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/bizflycloud/gobizfly"
	"github.com/bizflycloud/bizflycli/common"
)

// listCmd represents the list command
var volumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all volumes in your account",
	Long: `List all volumes in your BizFly Cloud account
Example: bizfly volume list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := apiClientForContext(cmd)
		volumes, err := client.Volume.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, vol := range volumes {
			data = append(data, []string{
				vol.ID, vol.Name, vol.Status, strconv.Itoa(vol.Size), vol.CreatedAt, vol.SnapshotID})
		}
		common.Output(volumeHeaderList, data)
	},
}

func init() {
	volumeCmd.AddCommand(volumeListCmd)
}
