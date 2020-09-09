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

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	firewallListHeader = []string{"ID", "Name", "Description", "Rules Count", "Servers Count", "Created at"}
	firewallAppliedServersHeader = []string{"ID", "Name", "Firewall ID"}
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "BizFly Cloud Firewall Interaction",
	Long:  "BizFly Cloud Firewall Action: Create, List, Delete, Update, Remove Server from Firewall",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("firewall called")
	},
}

var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all firewalls",
	Long: `List all firewalls of your account in a region
Example: bizfly firewall list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		firewalls, err := client.Firewall.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, firewall := range firewalls {
			fw := []string{firewall.ID, firewall.Name, firewall.Description, strconv.Itoa(firewall.RulesCount), strconv.Itoa(firewall.ServersCount), firewall.CreatedAt}
			data = append(data, fw)
		}
		formatter.Output(firewallListHeader, data)
	},
}

var firewallDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one firewall or more firewalls",
	Long: `Delete one firewall or more firewalls
Example: bizfly firewall delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58

You can delete multiple firewalls with list of firewall id
Example: bizfly firewall delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58 f5869e9c-9ab2-11ea-b9e3-e353a4f04836
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		for _, fwID := range args {
			fmt.Printf("Deleting firewall %s \n", fwID)
			_, err := client.Firewall.Delete(ctx, fwID)
			if err != nil {
				if errors.Is(err, gobizfly.ErrNotFound) {
					fmt.Printf("Firewall %s is not found", serverID)
					return
				} else {
					log.Fatal(err)
				}
			}
		}
	},
}

var firewallServerCmd = &cobra.Command{Use: "server"}

var firewallServerList = &cobra.Command{
	Use:   "list",
	Short: "List applied servers with the firewall",
	Long: `List applied servers with the firewall
Example: bizfly firewall server list  02b28284-5a18-4a0e-9ecc-d5d1acaf7e7b
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("You need to specify firewall ID in the command")
		}
		client, ctx := getApiClient(cmd)
		firewall, err := client.Firewall.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Firewall %s is not found", serverID)
				return
			} else {
				log.Fatal(err)
			}
		}
		var data [][]string
		for _, server := range firewall.Servers {
			fw := []string{server.ID, server.Name, firewall.ID}
			data = append(data, fw)
		}
		formatter.Output(firewallAppliedServersHeader, data)
	},
}

var firewallServerRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove server from a firewall",
	Long: `Remove server from a firewall
Example: bizfly firewall server remove <firewall ID> <server ID 1> <server ID 2> ..
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		if len(args) < 2 {
			log.Fatal("You need to specify firewall ID and server ID in the command")
		}
		client, ctx := getApiClient(cmd)
		frsr := gobizfly.FirewallRemoveServerRequest{
			Servers: args[1:],
		}
		fmt.Println(frsr)
		_, err := client.Firewall.RemoveServer(ctx, args[0], &frsr)
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Firewall %s is not found", serverID)
				return
			} else {
				log.Fatal(err)
			}
		}
		fmt.Println("Removed servers from a fitirewall completed")
		return
	},
}

func init() {
	rootCmd.AddCommand(firewallCmd)

	firewallCmd.AddCommand(firewallListCmd)
	firewallCmd.AddCommand(firewallDeleteCmd)
	firewallCmd.AddCommand(firewallServerCmd)

	firewallServerCmd.AddCommand(firewallServerRemove)
	firewallServerCmd.AddCommand(firewallServerList)
}
