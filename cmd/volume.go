/*
Copyright © (2020-2021) Bizfly Cloud

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
	"os"
	"strconv"
	"strings"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	volumeHeaderList  = []string{"ID", "Name", "Description", "Status", "Size", "Created At", "Type", "Snapshot ID", "Billing Plan"}
	volumeName        string
	volumeSize        int
	volumeType        string
	volumeCategory    string
	volumeBillingPlan string
	serverID          string
)

// volumeCmd represents the volume command
var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Bizfly Cloud Volume Interaction",
	Long:  `Bizfly Cloud Volume Action: Create, List, Delete, Extend Volume`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("volume called")
	},
}

// volumeDeleteCmd represents the delete command
var volumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete volume",
	Long: `Delete volume with volume ID as input
Example: bizfly volume delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58

You can delete multiple volumes with list of volume ID
Example: bizfly volume delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58 f5869e9c-9ab2-11ea-b9e3-e353a4f04836`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
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

// volumeGetCmd represents the get command
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
		client, ctx := getApiClient(cmd)

		volume, err := client.Volume.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Volume %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{volume.ID, volume.Name, volume.Description, volume.Status,
			strconv.Itoa(volume.Size), volume.CreatedAt, volume.SnapshotID, volume.BillingPlan})
		formatter.Output(volumeHeaderList, data)
	},
}

// volumeListCmd represents the list command
var volumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all volumes in your account",
	Long: `List all volumes in your Bizfly Cloud account
Example: bizfly volume list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		volumes, err := client.Volume.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, vol := range volumes {
			data = append(data, []string{
				vol.ID, vol.Name, vol.Description, vol.Status, strconv.Itoa(vol.Size),
				vol.CreatedAt, vol.VolumeType, vol.SnapshotID, vol.BillingPlan})
		}
		formatter.Output(volumeHeaderList, data)
	},
}

// volumeCreateCmd represents the create command
var volumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new volume",
	Long: `
Create a new volume
Use: bizfly volume create
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		vcr := gobizfly.VolumeCreateRequest{
			Name:             volumeName,
			Size:             volumeSize,
			VolumeType:       volumeType,
			SnapshotID:       snapshotID,
			ServerID:         serverID,
			AvailabilityZone: availabilityZone,
			VolumeCategory:   volumeCategory,
			Description:      description,
			BillingPlan:      volumeBillingPlan,
		}
		volume, err := client.Volume.Create(ctx, &vcr)
		if err != nil {
			fmt.Printf("Create a new volume error: %v", err)
			os.Exit(1)
		}
		var data [][]string
		data = append(data, []string{volume.ID, volume.Name, volume.Description, volume.Status,
			strconv.Itoa(volume.Size), volume.CreatedAt, volume.SnapshotID, volume.BillingPlan})
		formatter.Output(volumeHeaderList, data)
	},
}

// volumeAttachCmd represents the volume attach command
var volumeAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a volume to a server",
	Long: `
Attach a volume to a server
Use: bizfly volume attach <volume-id> <server-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Command error: use bizfly volume attach <volume-id> <server-id>")
			os.Exit(1)
		}
		volumeID := args[0]
		if volumeID == "" {
			fmt.Println("You need to specify volume-id in the command")
			os.Exit(1)
		}
		serverID := args[1]
		if serverID == "" {
			fmt.Println("You need to specify server-id in the command")
			os.Exit(1)
		}
		client, ctx := getApiClient(cmd)
		res, err := client.Volume.Attach(ctx, volumeID, serverID)
		if err != nil {
			fmt.Printf("Attach a volume to a server error: %v", err)
			os.Exit(1)
		}
		fmt.Println(res.Message)
	},
}

// volumeDetachCmd represents the volume detach command
var volumeDetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a volume from a server",
	Long: `
Detach a volume from a server
Use: bizfly volume detach <volume-id> <server-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Command error: use bizfly volume attach <volume-id> <server-id>")
			os.Exit(1)
		}
		volumeID := args[0]
		if volumeID == "" {
			fmt.Println("You need to specify volume-id in the command")
			os.Exit(1)
		}
		serverID := args[1]
		if serverID == "" {
			fmt.Println("You need to specify server-id in the command")
			os.Exit(1)
		}
		client, ctx := getApiClient(cmd)
		res, err := client.Volume.Detach(ctx, volumeID, serverID)
		if err != nil {
			fmt.Printf("Detach a volume from a server error: %v", err)
			os.Exit(1)
		}
		fmt.Println(res.Message)
	},
}

// extendVolumeCmd represent the resize volume command
var extendVolumeCmd = &cobra.Command{
	Use:   "extend",
	Short: "Extend size of a volume",
	Long: `
Extend size of a volume
Use: bizfly volume extend <volume-id> --size <new-size>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify the volume-id in the command. Use: bizfly volume extend <volume-id> --size <new size>")
			os.Exit(1)
		}
		volumeID := args[0]
		client, ctx := getApiClient(cmd)
		_, err := client.Volume.ExtendVolume(ctx, volumeID, volumeSize)
		if err != nil {
			fmt.Printf("Extend volume error: %v\n", err)
		}
		fmt.Printf("Extending volume %v\n", volumeID)
	},
}

var restoreVolumeCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore volume by using its snapshot",
	Long: `
Restore volume by using its snapshot
Use: bizfly volume restore <volume-id> --snapshot-id <snapshot-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify the volume-id in the command. Use: bizfly volume restore <volume-id> --snapshot-id <snapshot-id>")
			os.Exit(1)
		}
		volumeID := args[0]
		client, ctx := getApiClient(cmd)
		_, err := client.Volume.Restore(ctx, volumeID, snapshotID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Restoring volume %s using snapshot %s", volumeID, snapshotID)
	},
}

var patchVolumeCmd = &cobra.Command{
	Use:   "patch",
	Short: "Patch Volume",
	Long: `
Patch volume
Use: bizfly volume patch <volume-id> [--name <vol_name>] [--description <description>]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("You need to specify the volume-id in the command. Use: bizfly volume patch <volume-id> [--name <vol_name>] [--description <description>]")
		}
		volumeID := args[0]
		client, ctx := getApiClient(cmd)
		req := &gobizfly.VolumePatchRequest{}
		req.Description = description
		volume, err := client.Volume.Patch(ctx, volumeID, req)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{volume.ID, volume.Name, volume.Description, volume.Status, strconv.Itoa(volume.Size), volume.CreatedAt, volume.SnapshotID})
		formatter.Output(volumeHeaderList, data)
	},
}

func init() {
	rootCmd.AddCommand(volumeCmd)
	volumeCmd.AddCommand(volumeListCmd)
	volumeCmd.AddCommand(volumeGetCmd)
	volumeCmd.AddCommand(volumeDeleteCmd)

	vspf := restoreVolumeCmd.PersistentFlags()
	vspf.StringVar(&snapshotID, "snapshot-id", "", "Restore volume by using snapshot")
	_ = cobra.MarkFlagRequired(vspf, "snapshot-id")
	volumeCmd.AddCommand(restoreVolumeCmd)
	vcpf := volumeCreateCmd.PersistentFlags()
	vcpf.StringVar(&volumeName, "name", "", "Volume name")
	_ = cobra.MarkFlagRequired(vcpf, "name")
	vcpf.StringVar(&description, "description", "", "Volume description")
	vcpf.StringVar(&volumeType, "type", "HDD", "Volume type: SSD or HDD.")
	vcpf.StringVar(&volumeCategory, "category", "premium", "Volume category: premium, enterprise or basic.")
	vcpf.IntVar(&volumeSize, "size", 0, "Volume size")
	_ = cobra.MarkFlagRequired(vcpf, "size")
	vcpf.StringVar(&availabilityZone, "availability-zone", "HN1", "Avaialability Zone of volume.")
	vcpf.StringVar(&snapshotID, "snapshot-id", "", "Create a volume from a snapshot")
	vcpf.StringVar(&serverID, "server-id", "", "Create a new volume and attach to a server")
	vcpf.StringVar(&volumeBillingPlan, "billing-plan", "saving_plan", "Billing plan of volume: saving_plan, on_demand")
	volumeCmd.AddCommand(volumeCreateCmd)

	volumeCmd.AddCommand(volumeAttachCmd)

	volumeCmd.AddCommand(volumeDetachCmd)

	extendVolumeCmd.PersistentFlags().IntVar(&volumeSize, "size", 0, "Volume size")
	_ = cobra.MarkFlagRequired(extendVolumeCmd.PersistentFlags(), "size")
	volumeCmd.AddCommand(extendVolumeCmd)
	pvpf := patchVolumeCmd.PersistentFlags()
	pvpf.StringVar(&description, "description", "", "Patched volume description")
	_ = cobra.MarkFlagRequired(pvpf, "description")
	volumeCmd.AddCommand(patchVolumeCmd)
}
