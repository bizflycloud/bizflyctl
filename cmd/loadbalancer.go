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

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	lbListHeader       = []string{"ID", "Name", "Network Type", "IP Address", "Operating Status", "Type"}
	poolListHeader     = []string{"ID", "Name", "Algorithm", "Protocol", "Operating Status"}
	listenerListHeader = []string{"ID", "Name", "Protocol", "Protocol Port", "Operating Status", "Default Pool ID"}
)

// serverCmd represents the server command
var lbCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "BizFly Cloud Load Balancer Interaction",
	Long:  `BizFly Cloud Load Balancer Action: Create, List, Delete`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("loadbalancer called")
	},
}

var lbListenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "BizFly Cloud Load Balancer Listener Interaction",
	Long:  "BizFly Cloud Load Balancer Listener Action: Create, List, Delete, Get",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var lbPoolCmd = &cobra.Command{
	Use:   "pool",
	Short: "BizFly Cloud Load Balancer Pool Interaction",
	Long:  "BizFly Cloud Load Balancer Pool Action: Create, List, Delete, Get",
	Run:   func(cmd *cobra.Command, args []string) {},
}

// deleteCmd represents the delete command
var lbDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete load balancer",
	Long: `Delete Load Balancer with Load Balancer ID as input
Example: bizfly loadbalancer delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58

You can delete multiple loadbalancers with list of loadbalancer id
Example: bizfly loadbalancer delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58 f5869e9c-9ab2-11ea-b9e3-e353a4f04836
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		for _, lbID := range args {
			fmt.Printf("Deleting load balancer %s \n", lbID)
			lbdr := gobizfly.LoadBalancerDeleteRequest{ID: lbID, Cascade: true}
			err := client.LoadBalancer.Delete(ctx, &lbdr)
			if err != nil {
				if errors.Is(err, gobizfly.ErrNotFound) {
					fmt.Printf("Load Balancer %s is not found", lbID)
					return
				}
			}
		}
	},
}

// lbListCmd represents the list command
var lbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all load balancer in your account",
	Long:  `List all load balancer in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		lbs, err := client.LoadBalancer.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, lb := range lbs {
			s := []string{lb.ID, lb.Name, lb.NetworkType, lb.VipAddress, lb.OperatingStatus, lb.Type}
			data = append(data, s)
		}
		formatter.Output(lbListHeader, data)
	},
}

// lbGetCmd represents the get command
var lbGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a load balancer",
	Long: `Get detail a load balancer with load balancer ID as input
Example: bizfly loadbalancer get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)

		lb, err := client.LoadBalancer.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Load Balancer %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{lb.ID, lb.Name, lb.NetworkType, lb.VipAddress, lb.OperatingStatus, lb.Type})
		formatter.Output(lbListHeader, data)
	},
}

// lbPoolDeleteCmd represents the delete command
var lbPoolDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete pool in load balancer",
	Long: `Delete Pool in a Load balancer with Pool ID as input
Example: bizfly loadbalancer pool delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		// TODO: check length of args
		poolID := args[0]
		fmt.Printf("Deleting pool %s \n", poolID)
		err := client.Pool.Delete(ctx, poolID)
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Pool %s is not found", poolID)
				return
			}
		}
	},
}

// lbPoolListCmd represents the list command
var lbPoolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pools in a load balancer",
	Long: `List all pools in a load balancer
Example: bizfly loadbalancer pool list <loadbalancer_id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		// TODO Check length args
		lbID := args[0]
		pools, err := client.Pool.List(ctx, lbID, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, pool := range pools {
			s := []string{pool.ID, pool.Name, pool.LBAlgorithm, pool.Protocol, pool.OperatingStatus}
			data = append(data, s)
		}
		formatter.Output(poolListHeader, data)
	},
}

// lbPoolGetCmd represents the get command
var lbPoolGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get detail a pool in a load balancer",
	Long: `Get detail a pool in a load balancer with pool ID as input
Example: bizfly loadbalancer pool get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)

		pool, err := client.Pool.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Pool %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{pool.ID, pool.Name, pool.LBAlgorithm, pool.Protocol, pool.OperatingStatus})
		formatter.Output(poolListHeader, data)
	},
}

// lbListenerDeleteCmd represents the delete command
var lbListenerDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete listener in load balancer",
	Long: `Delete Listener in a Load balancer with Listener ID as input
Example: bizfly loadbalancer listener delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		// TODO: check length of args
		listenerID := args[0]
		fmt.Printf("Deleting listener %s \n", listenerID)
		err := client.Listener.Delete(ctx, listenerID)
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Listener %s is not found", listenerID)
				return
			}
		}
	},
}

// lbListenerListCmd represents the list command
var lbListenerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all listeners in a load balancer",
	Long: `List all listeners in a loadbalancer
Example: bizfly loadbalancer listener list <loadbalancer_id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		// TODO Check length args
		lbID := args[0]
		listeners, err := client.Listener.List(ctx, lbID, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, listener := range listeners {
			s := []string{listener.ID, listener.Name, listener.Protocol, strconv.Itoa(listener.ProtocolPort), listener.OperatingStatus, listener.DefaultPoolID}
			data = append(data, s)
		}
		formatter.Output(listenerListHeader, data)
	},
}

// lbListenerGetCmd represents the get command
var lbListenerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a listener in load balancer",
	Long: `Get detail a listener with listener  ID as input
Example: bizfly loadbalancer listener get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)

		listener, err := client.Listener.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Listener %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{listener.ID, listener.Name, listener.Protocol, strconv.Itoa(listener.ProtocolPort), listener.OperatingStatus, listener.DefaultPoolID})
		formatter.Output(listenerListHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(lbCmd)
	lbCmd.AddCommand(lbListCmd)
	lbCmd.AddCommand(lbGetCmd)
	lbCmd.AddCommand(lbDeleteCmd)

	lbCmd.AddCommand(lbPoolCmd)
	lbCmd.AddCommand(lbListenerCmd)

	lbPoolCmd.AddCommand(lbPoolGetCmd)
	lbPoolCmd.AddCommand(lbPoolListCmd)
	lbPoolCmd.AddCommand(lbPoolDeleteCmd)

	lbListenerCmd.AddCommand(lbListenerGetCmd)
	lbListenerCmd.AddCommand(lbListenerDeleteCmd)
	lbListenerCmd.AddCommand(lbListenerListCmd)
}
