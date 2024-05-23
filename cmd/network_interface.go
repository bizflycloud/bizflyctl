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
	"fmt"
	"log"
	"strings"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	networkInterfaceHeaders = []string{"ID", "Name", "Status", "Network ID", "Device ID", "IP Address",
		"IP Version", "Security Groups", "Created At", "Updated At"}
	vpcNetworkId           string
	networkInterfaceStatus string
	networkInterfaceType   string
	networkInterfaceName   string
	fixedIPAddress         string
	attachedServer         string
	firewallIDs            []string
)

var networkInterfaceCmd = &cobra.Command{
	Use:   "network-interface",
	Short: "Bizfly Cloud Network interfaces Interaction",
	Long:  `Bizfly Cloud Network interfaces Interaction: Create , List, Delete, Update, Action`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Network Interface called")
	},
}

var networkInterfaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Network Interfaces",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		opts := gobizfly.ListNetworkInterfaceOptions{
			VPCNetworkID: vpcNetworkId,
			Status:       networkInterfaceStatus,
			Type:         networkInterfaceType,
		}
		networkInterfaces, err := client.CloudServer.NetworkInterfaces().List(ctx, &opts)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		for _, networkInterface := range networkInterfaces {
			data = append(data, []string{
				networkInterface.ID,
				networkInterface.Name,
				networkInterface.Status,
				networkInterface.NetworkID,
				networkInterface.DeviceID,
				networkInterface.FixedIps[0].IPAddress,
				fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
				strings.Join(networkInterface.SecurityGroups, ","),
				networkInterface.CreatedAt,
				networkInterface.UpdatedAt,
			})
		}
		formatter.Output(networkInterfaceHeaders, data)
	},
}

var networkInterfaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Network Interface",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args[0]) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		payload := gobizfly.CreateNetworkInterfacePayload{
			Name:           networkInterfaceName,
			FixedIP:        fixedIPAddress,
			AttachedServer: attachedServer,
		}
		networkInterface, err := client.CloudServer.NetworkInterfaces().Create(ctx, args[0], &payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, []string{
			networkInterface.ID,
			networkInterface.Name,
			networkInterface.Status,
			networkInterface.NetworkID,
			networkInterface.DeviceID,
			networkInterface.FixedIps[0].IPAddress,
			fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
			strings.Join(networkInterface.SecurityGroups, ","),
			networkInterface.CreatedAt,
			networkInterface.UpdatedAt,
		})
		formatter.Output(networkInterfaceHeaders, data)
	},
}

var networkInterfaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Network Interface",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		networkInterface, err := client.CloudServer.NetworkInterfaces().Get(ctx, args[0])
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, []string{
			networkInterface.ID,
			networkInterface.Name,
			networkInterface.Status,
			networkInterface.NetworkID,
			networkInterface.DeviceID,
			networkInterface.FixedIps[0].IPAddress,
			fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
			strings.Join(networkInterface.SecurityGroups, ","),
			networkInterface.CreatedAt,
			networkInterface.UpdatedAt,
		})
		formatter.Output(networkInterfaceHeaders, data)
	},
}

var networkInterfaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Network Interface",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		err := client.CloudServer.NetworkInterfaces().Delete(ctx, args[0])
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("The Network Interface deleted successfully")
	},
}

var networkInterfaceAddFirewall = &cobra.Command{
	Use:   "add-firewalls",
	Short: "Add Firewall to the Network Interface",
	Long: `Add Firewall to the Network Interface: 
./bizfly network-interface add-firewalls <network-interface-id> --firewall <firewall-id> --firewall <firewall-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		payload := gobizfly.ActionNetworkInterfacePayload{
			Action:         "add_firewall",
			SecurityGroups: firewallIDs,
		}
		networkInterface, err := client.CloudServer.NetworkInterfaces().Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, []string{
			networkInterface.ID,
			networkInterface.Name,
			networkInterface.Status,
			networkInterface.NetworkID,
			networkInterface.DeviceID,
			networkInterface.FixedIps[0].IPAddress,
			fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
			strings.Join(networkInterface.SecurityGroups, ","),
			networkInterface.CreatedAt,
			networkInterface.UpdatedAt,
		})
		formatter.Output(networkInterfaceHeaders, data)
	},
}

var networkInterfaceRemoveFirewall = &cobra.Command{
	Use:   "remove-firewalls",
	Short: "Remove Firewall from the Network Interface",
	Long: `Remove Firewall from the Network Interface: 
./bizfly network-interface remove-firewalls <network-interface-id> --firewall <firewall-id> --firewall <firewall_id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		payload := gobizfly.ActionNetworkInterfacePayload{
			Action:         "remove_firewall",
			SecurityGroups: firewallIDs,
		}
		networkInterface, err := client.CloudServer.NetworkInterfaces().Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, []string{
			networkInterface.ID,
			networkInterface.Name,
			networkInterface.Status,
			networkInterface.NetworkID,
			networkInterface.DeviceID,
			networkInterface.FixedIps[0].IPAddress,
			fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
			strings.Join(networkInterface.SecurityGroups, ","),
			networkInterface.CreatedAt,
			networkInterface.UpdatedAt,
		})
		formatter.Output(networkInterfaceHeaders, data)
	},
}

var networkInterfaceAttachServer = &cobra.Command{
	Use:   "attach-server",
	Short: "Attach Server to the Network Interface",
	Long: `Attach Server to the Network Interface: 
./bizfly network-interface attach-server <network-interface-id> --server-id <server_id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		payload := gobizfly.ActionNetworkInterfacePayload{
			Action:   "attach_server",
			ServerID: serverID,
		}
		networkInterface, err := client.CloudServer.NetworkInterfaces().Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, []string{
			networkInterface.ID,
			networkInterface.Name,
			networkInterface.Status,
			networkInterface.NetworkID,
			networkInterface.DeviceID,
			networkInterface.FixedIps[0].IPAddress,
			fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
			strings.Join(networkInterface.SecurityGroups, ","),
			networkInterface.CreatedAt,
			networkInterface.UpdatedAt,
		})
		formatter.Output(networkInterfaceHeaders, data)
	},
}

var networkInterfaceDetachServer = &cobra.Command{
	Use:   "detach-server",
	Short: "Detach Server from the Network Interface",
	Long: `Detach Server from the Network Interface:
./bizfly network-interface detach-server <network-interface-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) == 0 {
			log.Fatal("Please provide Network Interface ID")
		}
		payload := gobizfly.ActionNetworkInterfacePayload{
			Action: "detach_server",
		}
		networkInterface, err := client.CloudServer.NetworkInterfaces().Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, []string{
			networkInterface.ID,
			networkInterface.Name,
			networkInterface.Status,
			networkInterface.NetworkID,
			networkInterface.DeviceID,
			networkInterface.FixedIps[0].IPAddress,
			fmt.Sprintf("%d", networkInterface.FixedIps[0].IPVersion),
			strings.Join(networkInterface.SecurityGroups, ","),
			networkInterface.CreatedAt,
			networkInterface.UpdatedAt,
		})
		formatter.Output(networkInterfaceHeaders, data)
	},
}

func init() {
	rootCmd.AddCommand(networkInterfaceCmd)
	networkInterfaceCmd.AddCommand(networkInterfaceGetCmd)
	networkInterfaceCmd.AddCommand(networkInterfaceDeleteCmd)
	networkInterfaceCmd.AddCommand(networkInterfaceListCmd)
	nilpf := networkInterfaceListCmd.PersistentFlags()
	nilpf.StringVar(&vpcNetworkId, "vpc-network-id", "", "VPC Network ID")
	nilpf.StringVar(&networkInterfaceStatus, "status", "", "Network Interface Status")
	nilpf.StringVar(&networkInterfaceType, "type", "", "Network Interface Type")

	networkInterfaceCmd.AddCommand(networkInterfaceCreateCmd)
	nicpf := networkInterfaceCreateCmd.PersistentFlags()
	nicpf.StringVar(&networkInterfaceName, "name", "", "Network Interface Name")
	nicpf.StringVar(&fixedIPAddress, "fixed-ip", "", "Fixed IP Address")
	nicpf.StringVar(&attachedServer, "server-id", "", "Attached Server ID")

	networkInterfaceCmd.AddCommand(networkInterfaceAddFirewall)
	niafpf := networkInterfaceAddFirewall.PersistentFlags()
	niafpf.StringArrayVar(&firewallIDs, "firewall-id", []string{}, "Firewall ID")

	networkInterfaceCmd.AddCommand(networkInterfaceRemoveFirewall)
	nirfpf := networkInterfaceRemoveFirewall.PersistentFlags()
	nirfpf.StringArrayVar(&firewallIDs, "firewall-id", []string{}, "Firewall ID")

	networkInterfaceCmd.AddCommand(networkInterfaceAttachServer)
	networkInterfaceAttachServer.PersistentFlags().StringVar(&serverID, "server-id", "", "Server ID")

	networkInterfaceCmd.AddCommand(networkInterfaceDetachServer)

}
