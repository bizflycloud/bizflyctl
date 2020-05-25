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
	"strings"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
)

var (
	lbListHeader = []string{"ID", "Name", "Network Type", "IP Address", "Operating Status", "Type"}
)

// serverCmd represents the server command
var lbCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "BizFly Cloud Load Balancer Interaction",
	Long:  `BizFly Cloud Load Balancer Action: Create, List, Delete`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

// deleteCmd represents the delete command
var lbDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Server",
	Long: `Delete Server with Load Balancer ID as input
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

func init() {
	rootCmd.AddCommand(lbCmd)
	lbCmd.AddCommand(lbListCmd)
	lbCmd.AddCommand(lbGetCmd)
	lbCmd.AddCommand(lbDeleteCmd)
}
