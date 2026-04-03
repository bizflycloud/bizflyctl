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
	internetGatewayHeaders = []string{"ID", "NAME", "STATUS", "VPC NAMES", "AVAILABILITY ZONES", "CREATED AT"}

	igwName             string
	igwNetworkIDs       []string
	igwDescription      string
	igwAvailabilityZone string
	igwNameFilter       string
)

var internetGatewayCmd = &cobra.Command{
	Use:   "internet-gateway",
	Short: "Bizfly Cloud Internet Gateway Interaction",
	Long:  `Bizfly Cloud Internet Gateway Interaction: Create, List, Get, Update, Delete`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Internet Gateway called")
	},
}

func parseIGWResult(igw *gobizfly.ExtendedInternetGateway) []string {
	vpcNames := []string{}
	for _, network := range igw.InterfacesInfo {
		if network.NetworkInfo == nil {
			continue
		}
		vpcNames = append(vpcNames, network.NetworkInfo.Name)
	}
	return []string{
		igw.ID,
		igw.Name,
		igw.Status,
		strings.Join(vpcNames, ","),
		strings.Join(igw.AvailabilityZoneHints, ","),
		igw.CreatedAt,
	}
}

var internetGatewayListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Internet Gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		detailed := true
		opts := gobizfly.ListInternetGatewayOpts{
			Detailed: &detailed,
		}
		if igwNameFilter != "" {
			opts.Name = &igwNameFilter
		}
		igws, err := client.CloudServer.InternetGateways().List(ctx, opts)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		for _, igw := range igws.InternetGateways {
			data = append(data, parseIGWResult(igw))
		}
		formatter.Output(internetGatewayHeaders, data)
	},
}

var internetGatewayGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide Internet Gateway ID")
		}
		client, ctx := getApiClient(cmd)
		igw, err := client.CloudServer.InternetGateways().Get(ctx, args[0])
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, parseIGWResult(igw))
		formatter.Output(internetGatewayHeaders, data)
	},
}

var internetGatewayCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if igwName == "" {
			log.Fatal("Please provide Internet Gateway name")
		}
		payload := gobizfly.CreateInternetGatewayPayload{
			Name:       igwName,
			NetworkIDs: &igwNetworkIDs,
		}
		if igwAvailabilityZone != "" {
			payload.AvailabilityZone = &igwAvailabilityZone
		}
		if igwDescription != "" {
			payload.Description = &igwDescription
		}
		igw, err := client.CloudServer.InternetGateways().Create(ctx, payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, parseIGWResult(igw))
		formatter.Output(internetGatewayHeaders, data)
	},
}

var internetGatewayUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide Internet Gateway ID")
		}
		client, ctx := getApiClient(cmd)
		igwID := args[0]
		oldIGW, err := client.CloudServer.InternetGateways().Get(ctx, igwID)
		if err != nil {
			log.Fatalln(err)
		}
		oldNetworkIDs := []string{}
		for _, network := range oldIGW.InterfacesInfo {
			if network.NetworkInfo == nil {
				continue
			}
			oldNetworkIDs = append(oldNetworkIDs, network.NetworkInfo.ID)
		}
		payload := gobizfly.UpdateInternetGatewayPayload{
			Name:        oldIGW.Name,
			Description: oldIGW.Description,
			NetworkIDs:  oldNetworkIDs,
		}
		if igwName != "" {
			payload.Name = igwName
		}
		if igwDescription != "" {
			payload.Description = igwDescription
		}
		if len(igwNetworkIDs) != 0 {
			payload.NetworkIDs = igwNetworkIDs
		}
		igw, err := client.CloudServer.InternetGateways().Update(ctx, args[0], payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, parseIGWResult(igw))
		formatter.Output(internetGatewayHeaders, data)
	},
}

var internetGatewayDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Internet Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide Internet Gateway ID")
		}
		client, ctx := getApiClient(cmd)
		err := client.CloudServer.InternetGateways().Delete(ctx, args[0])
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Internet Gateway deleted successfully")
	},
}

var internetGatewayDetachVPCCmd = &cobra.Command{
	Use:   "detach-vpc",
	Short: "Detach VPC out of Internet Gateway",
	Long: `Detach VPC out of Internet Gateway by setting network IDs to empty.
Usage: ./bizfly internet-gateway detach-vpc <internet-gateway-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide Internet Gateway ID")
		}
		client, ctx := getApiClient(cmd)
		igwID := args[0]
		oldIGW, err := client.CloudServer.InternetGateways().Get(ctx, igwID)
		if err != nil {
			log.Fatalln(err)
		}
		payload := gobizfly.UpdateInternetGatewayPayload{
			Name:        oldIGW.Name,
			Description: oldIGW.Description,
			NetworkIDs:  []string{},
		}
		igw, err := client.CloudServer.InternetGateways().Update(ctx, igwID, payload)
		if err != nil {
			log.Fatalln(err)
		}
		var data [][]string
		data = append(data, parseIGWResult(igw))
		formatter.Output(internetGatewayHeaders, data)
	},
}

func init() {
	rootCmd.AddCommand(internetGatewayCmd)
	internetGatewayCmd.AddCommand(internetGatewayListCmd)
	igwlf := internetGatewayListCmd.PersistentFlags()
	igwlf.StringVar(&igwNameFilter, "name", "", "Filter by name")

	internetGatewayCmd.AddCommand(internetGatewayGetCmd)
	internetGatewayCmd.AddCommand(internetGatewayDeleteCmd)
	internetGatewayCmd.AddCommand(internetGatewayDetachVPCCmd)

	internetGatewayCmd.AddCommand(internetGatewayCreateCmd)
	igwcpf := internetGatewayCreateCmd.PersistentFlags()
	igwcpf.StringVar(&igwName, "name", "", "The name of the Internet Gateway")
	igwcpf.StringArrayVar(&igwNetworkIDs, "network-id", []string{}, "The ID of VPC network")
	igwcpf.StringVar(&igwDescription, "description", "", "The description of the Internet Gateway")
	igwcpf.StringVar(&igwAvailabilityZone, "availability-zone", "", "The availability zone of the Internet Gateway")

	internetGatewayCmd.AddCommand(internetGatewayUpdateCmd)
	igwupf := internetGatewayUpdateCmd.PersistentFlags()
	igwupf.StringVar(&igwName, "name", "", "The name of the Internet Gateway")
	igwupf.StringArrayVar(&igwNetworkIDs, "network-id", []string{}, "The ID of VPC network")
	igwupf.StringVar(&igwDescription, "description", "", "Description")
}
