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
	"strconv"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	firewallListHeader           = []string{"ID", "Name", "Description", "Rules Count", "Servers Count", "Created at"}
	firewallAppliedServersHeader = []string{"ID", "Name", "Firewall ID"}
	firewallRuleHeader           = []string{"ID", "Description", "Direction", "Type", "Ether Type", "Protocol", "CIDR", "Port Range", "Remote IP Prefix"}

	fwRuleDirection string
	fwRuleProtocol  string
	fwRuleCIDR      string
	fwPortRange     string
	fwName          string
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Bizfly Cloud Firewall Interaction",
	Long:  "Bizfly Cloud Firewall Action: Create, List, Delete, Update, Remove Server from Firewall",
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
		firewalls, err := client.CloudServer.Firewalls().List(ctx, &gobizfly.ListOptions{})
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
			_, err := client.CloudServer.Firewalls().Delete(ctx, fwID)
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
		firewall, err := client.CloudServer.Firewalls().Get(ctx, args[0])
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
		_, err := client.CloudServer.Firewalls().RemoveServer(ctx, args[0], &frsr)
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Firewall %s is not found", serverID)
				return
			} else {
				log.Fatal(err)
			}
		}
		fmt.Println("Removed servers from a fitirewall completed")
	},
}

var firewallCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new firewall",
	Long: `Create a new firewall in your account
Example:
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		firewall, err := client.CloudServer.Firewalls().Create(ctx, &gobizfly.FirewallRequestPayload{Name: fwName})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		fw := []string{firewall.ID, firewall.Name, firewall.Description, strconv.Itoa(firewall.RulesCount), strconv.Itoa(firewall.ServersCount), firewall.CreatedAt}
		data = append(data, fw)
		formatter.Output(firewallListHeader, data)
	},
}

var firewallRuleCmd = &cobra.Command{
	Use: "rule",
}

var firewallRuleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rules in the firewall",
	Long: `List all rules in the firewall
Example: bizfly firewall rule list <firwall id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("You need to specify firewall ID in the command")
		}
		client, ctx := getApiClient(cmd)
		firewall, err := client.CloudServer.Firewalls().Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Firewall %s is not found", serverID)
				return
			} else {
				log.Fatal(err)
			}
		}
		var data [][]string
		for _, rule := range firewall.InBound {
			data = append(data, []string{rule.ID, rule.Description, rule.Direction, rule.Type, rule.EtherType, rule.Protocol, rule.CIDR, rule.PortRange, rule.RemoteIPPrefix})
		}
		for _, rule := range firewall.OutBound {
			data = append(data, []string{rule.ID, rule.Description, rule.Direction, rule.Type, rule.EtherType, rule.Protocol, rule.CIDR, rule.PortRange, rule.RemoteIPPrefix})
		}
		formatter.Output(firewallRuleHeader, data)
	},
}

var firewallRuleDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a rule in a firewall",
	Long: `Delete a rule in a firewall
Example: bizfly firewall rule delete <firewall ID> <rule ID>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("You need to specify firewall ID and rule ID in the command")
		}
		client, ctx := getApiClient(cmd)
		_, err := client.CloudServer.Firewalls().Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Firewall %s is not found", serverID)
				return
			} else {
				log.Fatal(err)
			}
		}
		resp, err := client.CloudServer.Firewalls().DeleteRule(ctx, args[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.Message)
	},
}

var firewallRuleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new firewall rule",
	Long: `Create a new rule in your firewall
Example: bizfly firewall rule create <firewall ID> --direction <ingress|egress> --protocol <tcp|udp> --port-range <port range> --cidr <CIDR>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("You need to specify the fireewall ID in the command")
		}
		client, ctx := getApiClient(cmd)
		frcr := gobizfly.FirewallSingleRuleCreateRequest{
			Direction: fwRuleDirection,
			FirewallRuleCreateRequest: gobizfly.FirewallRuleCreateRequest{
				Protocol: fwRuleProtocol,
				CIDR:     fwRuleCIDR,
				Type:     "CUSTOM",
			},
		}
		if fwPortRange != "" {
			frcr.PortRange = fwPortRange
		}
		resp, err := client.CloudServer.Firewalls().CreateRule(ctx, args[0], &gobizfly.FirewallSingleRuleCreateRequest{
			Direction: fwRuleDirection,
			FirewallRuleCreateRequest: gobizfly.FirewallRuleCreateRequest{
				Protocol:  fwRuleProtocol,
				CIDR:      fwRuleCIDR,
				PortRange: fwPortRange,
				Type:      "CUSTOM",
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created new firewall rule with ID %s", resp.ID)
	},
}

func init() {
	rootCmd.AddCommand(firewallCmd)

	firewallCmd.AddCommand(firewallListCmd)
	firewallCmd.AddCommand(firewallDeleteCmd)
	firewallCmd.AddCommand(firewallServerCmd)

	firewallServerCmd.AddCommand(firewallServerRemove)
	firewallServerCmd.AddCommand(firewallServerList)

	firewallCmd.AddCommand(firewallRuleCmd)
	firewallRuleCmd.AddCommand(firewallRuleListCmd)
	firewallRuleCmd.AddCommand(firewallRuleDeleteCmd)
	firewallRuleCmd.AddCommand(firewallRuleCreateCmd)
	frcf := firewallRuleCreateCmd.PersistentFlags()
	frcf.StringVar(&fwRuleDirection, "direction", "", "Direction, ingress or egress")
	_ = cobra.MarkFlagRequired(frcf, "direction")
	frcf.StringVar(&fwRuleProtocol, "protocol", "", "Protocol, one of tcp and udp")
	_ = cobra.MarkFlagRequired(frcf, "protocol")
	frcf.StringVar(&fwPortRange, "port-range", "1-65535", "Port or Port range. You can specify only one port or port range. Example: 80 and 80-90.")
	frcf.StringVar(&fwRuleCIDR, "cidr", "0.0.0.0/0", "CIDR. Example: 10.0.0.0/24")

	firewallCmd.AddCommand(firewallCreateCmd)
	fcf := firewallCreateCmd.PersistentFlags()
	fcf.StringVar(&fwName, "name", "", "Firewall name")
	_ = cobra.MarkFlagRequired(fcf, "name")

}
