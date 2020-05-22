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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
)

var (
	volumeHeaderList = []string{"ID", "Name", "Status", "Size", "Created At", "Type", "Snapshot ID"}
)

// volumeCmd represents the volume command
var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "BizFly Cloud Volume Interaction",
	Long: `BizFly Cloud Server Action: Create, List, Delete, Extend Volume`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("volume called")
	},
}

// deleteCmd represents the delete command
var volumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete volume",
	Long: `Delete volume with volume ID as input
Example: bizfly volume delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58

You can delete multiple volumes with list of volume ID
Example: bizfly volume delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58 f5869e9c-9ab2-11ea-b9e3-e353a4f04836`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := apiClientForContext(cmd)
		for _, volumeID := range args {
			fmt.Printf("Deleting volume %s \n", volumeID)
			err := client.Volume.Delete(ctx, volumeID)
			if err != nil {
				if errors.Is(err, gobizfly.ErrNotFound) {
					fmt.Printf("Volume %s is not found", volumeID)
					return
				}
			}
		}
	},
}

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
		formatter.Output(volumeHeaderList, data)
	},
}

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
				vol.ID, vol.Name, vol.Status, strconv.Itoa(vol.Size), vol.CreatedAt, vol.VolumeType, vol.SnapshotID})
		}
		formatter.Output(volumeHeaderList, data)
	},
}

func init() {
	rootCmd.AddCommand(volumeCmd)
	volumeCmd.AddCommand(volumeListCmd)
	volumeCmd.AddCommand(volumeGetCmd)
	volumeCmd.AddCommand(volumeDeleteCmd)
}
