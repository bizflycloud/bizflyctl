/*
Copyright © 2021 BizFly Cloud

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
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var (
	wanIPHeader = []string{"Id", "Name", "Status", "Device Id", "IP Address", "IP Version", "Billing Type", "Bandwidth",
		"Zone", "Created At", "Updated At"}
	wanIpName string
)

var wanIPCmd = &cobra.Command{
	Use:   "wan-ip",
	Short: "Bizfly Cloud WAN IP Interaction",
	Long:  `Bizfly Cloud WAN IP Interaction: Create, Delete, List, Get, Action`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WAN IP called")
	},
}

var wanIpListCmd = &cobra.Command{
	Use:   "list",
	Short: "List WAN IP",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		wanIps, err := client.WanIP.List(ctx)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, wanIp := range wanIps {
			data = append(data, []string{
				wanIp.ID,
				wanIp.Name,
				wanIp.Status,
				wanIp.DeviceID,
				wanIp.IpAddress,
				strconv.Itoa(wanIp.IpVersion),
				wanIp.BillingType,
				strconv.Itoa(wanIp.Bandwidth),
				wanIp.AvailabilityZone,
				wanIp.CreatedAt,
				wanIp.UpdatedAt,
			})
		}
		formatter.Output(wanIPHeader, data)
	},
}

var wanIPCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create WAN IP",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		payload := gobizfly.CreateWanIpPayload{
			Name:             wanIpName,
			AvailabilityZone: availabilityZone,
			AttachedServer:   serverID,
		}

		wanIp, err := client.WanIP.Create(ctx, &payload)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			wanIp.ID,
			wanIp.Name,
			wanIp.Status,
			wanIp.DeviceID,
			wanIp.IpAddress,
			strconv.Itoa(wanIp.IpVersion),
			wanIp.BillingType,
			strconv.Itoa(wanIp.Bandwidth),
			wanIp.AvailabilityZone,
			wanIp.CreatedAt,
			wanIp.UpdatedAt,
		})
		formatter.Output(wanIPHeader, data)
	},
}

var wanIPGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get WAN IP",
	Long:  `Get WAN IP: ./bizfly wan-ip get <wan-ip-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		client, ctx := getApiClient(cmd)
		wanIp, err := client.WanIP.Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			wanIp.ID,
			wanIp.Name,
			wanIp.Status,
			wanIp.DeviceID,
			wanIp.IpAddress,
			strconv.Itoa(wanIp.IpVersion),
			wanIp.BillingType,
			strconv.Itoa(wanIp.Bandwidth),
			wanIp.AvailabilityZone,
			wanIp.CreatedAt,
			wanIp.UpdatedAt,
		})
		formatter.Output(wanIPHeader, data)
	},
}

var wanIpDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete WAN IP",
	Long:  `Delete WAN IP: ./bizfly wan-ip delete <wan-ip-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		client, ctx := getApiClient(cmd)
		err := client.WanIP.Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("The WAN IP is deleted")
	},
}

var wanIpAttachServerCmd = &cobra.Command{
	Use:   "attach-server",
	Short: "Attach WAN IP to server",
	Long:  `Attach WAN IP to server: ./bizfly wan-ip attach-server <wan-ip-id> --server-id <server-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		client, ctx := getApiClient(cmd)
		payload := gobizfly.ActionWanIpPayload{
			Action:   "attach_server",
			ServerId: serverID,
		}
		err := client.WanIP.Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatal(err)
		}
		wanIp, err := client.WanIP.Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			wanIp.ID,
			wanIp.Name,
			wanIp.Status,
			wanIp.DeviceID,
			wanIp.IpAddress,
			strconv.Itoa(wanIp.IpVersion),
			wanIp.BillingType,
			strconv.Itoa(wanIp.Bandwidth),
			wanIp.AvailabilityZone,
			wanIp.CreatedAt,
			wanIp.UpdatedAt,
		})
		formatter.Output(wanIPHeader, data)
	},
}

var wanIpDetachServerCmd = &cobra.Command{
	Use:   "detach-server",
	Short: "Detach the WAN IP from server",
	Long:  `Detach the WAN IP from server: ./bizfly wan-ip detach-server <wan-ip-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		client, ctx := getApiClient(cmd)
		payload := gobizfly.ActionWanIpPayload{
			Action: "detach_server",
		}
		err := client.WanIP.Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatal(err)
		}
		wanIp, err := client.WanIP.Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			wanIp.ID,
			wanIp.Name,
			wanIp.Status,
			wanIp.DeviceID,
			wanIp.IpAddress,
			strconv.Itoa(wanIp.IpVersion),
			wanIp.BillingType,
			strconv.Itoa(wanIp.Bandwidth),
			wanIp.AvailabilityZone,
			wanIp.CreatedAt,
			wanIp.UpdatedAt,
		})
		formatter.Output(wanIPHeader, data)
	},
}

var wanIpConvertToPaidCmd = &cobra.Command{
	Use:   "convert-to-paid",
	Short: "Convert WAN IP to paid one",
	Long:  `Convert WAN IP to paid one: ./bizfly wan-ip convert-to-paid <wan-ip-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		client, ctx := getApiClient(cmd)
		payload := gobizfly.ActionWanIpPayload{
			Action: "convert_to_paid",
		}
		err := client.WanIP.Action(ctx, args[0], &payload)
		if err != nil {
			log.Fatal(err)
		}
		wanIp, err := client.WanIP.Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{
			wanIp.ID,
			wanIp.Name,
			wanIp.Status,
			wanIp.DeviceID,
			wanIp.IpAddress,
			strconv.Itoa(wanIp.IpVersion),
			wanIp.BillingType,
			strconv.Itoa(wanIp.Bandwidth),
			wanIp.AvailabilityZone,
			wanIp.CreatedAt,
			wanIp.UpdatedAt,
		})
		formatter.Output(wanIPHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(wanIPCmd)
	wanIPCmd.AddCommand(wanIpListCmd)
	wanIPCmd.AddCommand(wanIPGetCmd)
	wanIPCmd.AddCommand(wanIpDeleteCmd)

	wicpf := wanIPCreateCmd.PersistentFlags()
	wicpf.StringVar(&availabilityZone, "zone", "", "The availability zone")
	wicpf.StringVar(&wanIpName, "name", "", "The WAN IP name")
	wicpf.StringVar(&serverID, "server-id", "", "The server id want to attach")
	wanIPCmd.AddCommand(wanIPCreateCmd)

	wiaspf := wanIpAttachServerCmd.PersistentFlags()
	wiaspf.StringVar(&serverID, "server-id", "", "The server id want to attach")
	wanIPCmd.AddCommand(wanIpAttachServerCmd)
	wanIPCmd.AddCommand(wanIpDetachServerCmd)
	wanIPCmd.AddCommand(wanIpConvertToPaidCmd)
}
