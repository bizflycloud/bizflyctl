/*
Copyright © (2021) Bizfly Cloud

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
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	vpcListHeader = []string{"ID", "Name", "MTU", "CIDR", "Description", "Tags", "Created At", "Is Default"}
	vpcName       string
	vpcID         string
	description   string
	cidr          string
	isDefault     bool
)

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Bizfly Virtual Private Network Interaction",
	Long:  "Bizfly Virtual Private Network Action: Create, List, Delete, Update",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vpc called")
	},
}

var vpcDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete VPC",
	Long: `Delete VPC with vpc ID as input
Example: bizfly vpc delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)

		fmt.Printf("Deleting VPC: %v", vpcID)
		err := client.VPC.Delete(ctx, vpcID)
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("VPC %s is not found", serverID)
			} else {
				fmt.Printf("Error when delete VPC %v", err)
				return
			}
		}
	},
}

var vpcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vpcs in your account",
	Long:  "List all vpcs in your account",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		vpcs, err := client.VPC.List(ctx)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, vpc := range vpcs {
			s := []string{vpc.ID, vpc.Name, strconv.Itoa(vpc.MTU), vpc.Subnets[0].CIDR, vpc.Description,
				strings.Join(vpc.Tags, ", "), vpc.CreatedAt, strconv.FormatBool(vpc.IsDefault)}
			data = append(data, s)
		}
		formatter.Output(vpcListHeader, data)
	},
}

var vpcGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a VPC",
	Long: `Get detail a VPC with VPC ID as input
Example: bizfly vpc get fd554aac-9ab1-11ea-b09d-bbaf82f02f58`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknown variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		vpc, err := client.VPC.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Server %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		s := []string{vpc.ID, vpc.Name, strconv.Itoa(vpc.MTU), vpc.Subnets[0].CIDR, vpc.Description,
			strings.Join(vpc.Tags, ", "), vpc.CreatedAt, strconv.FormatBool(vpc.IsDefault)}
		data = append(data, s)
		formatter.Output(vpcListHeader, data)
	},
}

var vpcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a VPC",
	Long:  "Create a new VPC, return its properties",
	Run: func(cmd *cobra.Command, args []string) {
		if vpcName == "" {
			fmt.Println("You need to specify VPC name to create a new VPC")
		}
		cvpl := gobizfly.CreateVPCPayload{
			Name:        vpcName,
			Description: description,
			CIDR:        cidr,
			IsDefault:   isDefault,
		}
		client, ctx := getApiClient(cmd)
		vpc, err := client.VPC.Create(ctx, &cvpl)
		if err != nil {
			fmt.Printf("Create VPC error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Create VPC successfully\n")
		var data [][]string
		s := []string{vpc.ID, vpc.Name, strconv.Itoa(vpc.MTU), vpc.Subnets[0].CIDR, vpc.Description,
			strings.Join(vpc.Tags, ", "), vpc.CreatedAt, strconv.FormatBool(vpc.IsDefault)}
		data = append(data, s)
		formatter.Output(vpcListHeader, data)
	},
}

var vpcUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a VPC",
	Long:  "Update a VPC",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify vpc-id in the command. Use bizfly vpc update <vpc-id> ...")
		}
		uvpl := gobizfly.UpdateVPCPayload{
			Name:        vpcName,
			Description: description,
			CIDR:        cidr,
			IsDefault:   isDefault,
		}
		client, ctx := getApiClient(cmd)
		vpc, err := client.VPC.Update(ctx, args[0], &uvpl)
		if err != nil {
			fmt.Printf("Update VPC error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Update VPC successfully\n")
		var data [][]string
		s := []string{vpc.ID, vpc.Name, strconv.Itoa(vpc.MTU), vpc.Subnets[0].CIDR, vpc.Description,
			strings.Join(vpc.Tags, ", "), vpc.CreatedAt, strconv.FormatBool(vpc.IsDefault)}
		data = append(data, s)
		formatter.Output(vpcListHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(vpcCmd)
	vpcCmd.AddCommand(vpcListCmd)
	vpcCmd.AddCommand(vpcGetCmd)
	vpcCmd.AddCommand(vpcDeleteCmd)

	vcpf := vpcCreateCmd.PersistentFlags()
	vcpf.StringVar(&vpcName, "name", "", "Name of VPC")
	_ = cobra.MarkFlagRequired(vcpf, "name")
	vcpf.StringVar(&description, "description", "", "Description")
	vcpf.StringVar(&cidr, "cidr", "", "CIDR")
	vcpf.BoolVar(&isDefault, "is-default", false, "Is default")
	vpcCmd.AddCommand(vpcCreateCmd)

	vupf := vpcUpdateCmd.PersistentFlags()
	vupf.StringVar(&vpcName, "name", "", "Name of VPC")
	vupf.StringVar(&description, "description", "", "Description")
	vupf.StringVar(&cidr, "cidr", "", "CIDR")
	vupf.BoolVar(&isDefault, "is-default", false, "Is default")
	vpcCmd.AddCommand(vpcUpdateCmd)

}
