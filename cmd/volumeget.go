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
	"github.com/bizflycloud/bizflycli/common"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

// getCmd represents the get command
var volumeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get detail a volume",
	Long: `Get detail a volume in your account
Example: bizfly volume get 9e580b1a-0526-460b-9a6f-d8f80130bda8
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := apiClientForContext(cmd)

		volume, err := client.Volume.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Volume %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{volume.ID, volume.Name, volume.Status, strconv.Itoa(volume.Size), volume.CreatedAt, volume.SnapshotID})
		common.Output(volumeHeaderList, data)
	},
}

func init() {
	volumeCmd.AddCommand(volumeGetCmd)

}
