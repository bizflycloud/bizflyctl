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
	"os"
	"strings"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	serverListHeader = []string{"ID", "Name", "Key Name", "Status", "Flavor", "LAN IP", "WAN IP", "Created at"}

	serverName string
	// serverOS gobizfly type

	imageID    string
	volumeID   string
	snapshotID string
	flavorName string

	// basic, premium, enterprise category
	serverCategory   string
	availabilityZone string

	// rootdisk
	rootDiskType string
	rootDiskSize int
	// ssh key
	sshKey string
)

//type responseMessage struct {
//	message string `json:"message"`
//}

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
		client, ctx := getApiClient(cmd)
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
		client, ctx := getApiClient(cmd)
		servers, err := client.Server.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, server := range servers {
			var LanIP []string
			for _, lan := range server.IPAddresses.LanAddresses {
				LanIP = append(LanIP, lan.Address)
			}
			LanIPAddrs := strings.Join(LanIP, ", ")
			var WanIP []string
			for _, wanv4 := range server.IPAddresses.WanV4Addresses {
				WanIP = append(WanIP, wanv4.Address)
			}
			for _, wanv6 := range server.IPAddresses.WanV6Addresses {
				WanIP = append(WanIP, wanv6.Address)
			}
			WanIPAddrs := strings.Join(WanIP, ", ")
			s := []string{server.ID, server.Name, server.KeyName, server.Status, server.Flavor.Name, LanIPAddrs, WanIPAddrs, server.CreatedAt}
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
		client, ctx := getApiClient(cmd)

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

// serverCreateCmd represents the create server command
var serverCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a server",
	Long:  "Create a new server, return a task ID of the processing",
	Run: func(cmd *cobra.Command, arg []string) {

		if imageID == "" && volumeID == "" && snapshotID == "" {
			fmt.Println("You need to specify image-id or volume-id or snapshot-id to create a new server")
		}

		var serverOS gobizfly.ServerOS

		if imageID != "" {
			serverOS.Type = "image"
			serverOS.ID = imageID
		}
		if volumeID != "" {
			serverOS.Type = "volume"
			serverOS.ID = volumeID
		}

		if snapshotID != "" {
			serverOS.Type = "snapshot"
			serverOS.ID = snapshotID
		}

		rootDisk := gobizfly.ServerDisk{
			Type: rootDiskType,
			Size: rootDiskSize,
		}

		scr := gobizfly.ServerCreateRequest{
			Name:             serverName,
			FlavorName:       flavorName,
			SSHKey:           sshKey,
			RootDisk:         &rootDisk,
			Type:             serverCategory,
			AvailabilityZone: availabilityZone,
			OS:               &serverOS,
		}
		client, ctx := getApiClient(cmd)
		svrTask, err := client.Server.Create(ctx, &scr)
		if err != nil {
			fmt.Printf("Create server error: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Creating server with task id: %v\n", svrTask.Task[0])
	},
}

// serverRebootCmd represents the reboot server command
var serverRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot a server. This is soft reboot",
	Long: `
Reboot a server
Use: bizfly server reboot <server-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify server-id in the command. Use bizfly server reboot <server-id>")
			os.Exit(1)
		}
		serverID := args[0]
		client, ctx := getApiClient(cmd)
		res, err := client.Server.SoftReboot(ctx, serverID)
		if err != nil {
			fmt.Printf("Reboot server error %v", err)
			os.Exit(1)
		}
		fmt.Println(res.Message)

	},
}

// serverHardRebootCmd represents the hard reboot server command
var serverHardRebootCmd = &cobra.Command{
	Use:   "hard reboot",
	Short: "Hard reboot a server",
	Long: `
Hard reboot a server.
Use: bizfly server hard reboot <server-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("You need to specify server-id in the command. Use bizfly server hard reboot <server-id>")
			os.Exit(1)
		}
		serverID := args[1]
		client, ctx := getApiClient(cmd)
		res, err := client.Server.HardReboot(ctx, serverID)
		if err != nil {
			fmt.Printf("Hard Reboot server error %v\n", err)
			os.Exit(1)
		}
		fmt.Println(res.Message)

	},
}

// serverStopCmd represents the hard stop server command
var serverStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a server",
	Long: `
Stop a server.
Use: bizfly server stop <server-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify server-id in the command. Use bizfly server stop <server-id>")
			os.Exit(1)
		}
		serverID := args[0]
		client, ctx := getApiClient(cmd)
		_, err := client.Server.Stop(ctx, serverID)
		if err != nil {
			fmt.Printf("Stop server error %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Stopping server: %s\n", serverID)

	},
}

// serverStartCmd represents the hard stop server command
var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a server",
	Long: `
Start a server.
Use: bizfly server start <server-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify server-id in the command. Use bizfly server start <server-id>")
			os.Exit(1)
		}
		serverID := args[0]
		client, ctx := getApiClient(cmd)
		_, err := client.Server.Start(ctx, serverID)
		if err != nil {
			fmt.Printf("Start server error %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Starting server: %s\n", serverID)

	},
}

// serverResizeCmd represents the hard stop server command
var serverResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a server",
	Long: `
Resize a server.
Use: bizfly server resize <server-id> --flavor <flavor name>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify server-id in the command. Use bizfly server resize <server-id> --flavor")
			os.Exit(1)
		}
		serverID := args[0]
		client, ctx := getApiClient(cmd)
		_, err := client.Server.Resize(ctx, serverID, flavorName)
		if err != nil {
			fmt.Printf("Resize server error %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Resizing server: %s\n", serverID)

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverGetCmd)
	serverCmd.AddCommand(serverDeleteCmd)

	scpf := serverCreateCmd.PersistentFlags()
	scpf.StringVar(&serverName, "name", "", "Name of server")
	_ = cobra.MarkFlagRequired(scpf, "name")
	scpf.StringVar(&imageID, "image-id", "", "ID of OS image. Create a root disk using this image ID")
	scpf.StringVar(&volumeID, "volume-id", "", "ID of volume. Create a server using an existing root disk volume.")
	scpf.StringVar(&snapshotID, "snapshot-id", "", "ID of snapshot. Create a server from a snapshot ID.")
	scpf.StringVar(&flavorName, "flavor", "", "Name of flavor. Flavor for create a server. Using 'bizfly flavor list' to get a list of flavors")
	_ = cobra.MarkFlagRequired(scpf, "flavor")
	scpf.StringVar(&serverCategory, "category", "premium", "Server category: basic, premium or enterprise.")
	scpf.StringVar(&availabilityZone, "availability-zone", "HN1", "Availability Zone of server.")
	scpf.StringVar(&rootDiskType, "rootdisk-type", "HDD", "Type of root disk: HDD or SSD.")
	scpf.IntVar(&rootDiskSize, "rootdisk-size", 0, "Size of root disk in Gigabyte. Minimum is 20GB")
	_ = cobra.MarkFlagRequired(scpf, "rootdisk-size")
	scpf.StringVar(&sshKey, "ssh-key", "", "SSH key")

	serverCmd.AddCommand(serverCreateCmd)
	serverCmd.AddCommand(serverRebootCmd)
	serverCmd.AddCommand(serverHardRebootCmd)
	serverCmd.AddCommand(serverStopCmd)
	serverCmd.AddCommand(serverStartCmd)

	serverResizeCmd.PersistentFlags().StringVar(&flavorName, "flavor", "", "Name of flavor.")
	_ = cobra.MarkFlagRequired(serverResizeCmd.PersistentFlags(), "flavor")
	serverCmd.AddCommand(serverResizeCmd)
}
