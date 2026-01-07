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
	kafkaListClusterHeader   = []string{"ID", "Name", "Zone", "Num Of Nodes", "Volume Size", "Public Access", "Status", "Flavor", "Created At"}
	kafkaDetailClusterHeader = []string{"ID", "Name", "Zone", "Flavor", "Disk Usage/Capacity", "LAN Connect", "WAN Connect", "Status", "Created At"}
	kafkaListFlavorHeader    = []string{"ID", "Name", "VCPUs", "RAM", "Disk Capacity", "Category"}
	kafkaListVersionHeader   = []string{"ID", "Name", "Code", "Default"}

	nodes                  int
	publicAccess           bool
	kafkaVolumeSize        int
	kafkaVersion           string
	kafkaFlavor            string
	kafkaResizeType        string
	kafkaNetworkInterfaces string
	kafkaAvailabilityZone  string
	kafkaSizeType          string
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Bizfly Cloud Kafka Interaction",
	Long:  `Bizfly Cloud Kafka Infomation: Cluster, Flavor, Version`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // Display the help message
	},
}

// kafkaClusterCmd represents the kafka command
var kafkaClusterCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Bizfly Cloud Kafka Interaction",
	Long:  `Bizfly Cloud Kafka Action: Cluster, Flavor, Version`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // Display the help message
	},
}

// kafkaClusterListCmd represents the list cluster command
var kafkaClusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cluster in your account",
	Long:  `List all cluster in your account
	Example: bizfly kafka clusters list`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		clusters, err := client.Kafka.List(ctx, &gobizfly.KafkaClusterListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, cluster := range clusters {
			data = append(data, []string{cluster.ID, cluster.Name, cluster.AvailabilityZone, strconv.Itoa(cluster.Nodes), fmt.Sprintf("%d GB", cluster.VolumeSize), strconv.FormatBool(cluster.PublicAccess), cluster.Status, cluster.Flavor, cluster.CreatedAt})
		}
		formatter.Output(kafkaListClusterHeader, data)
	},
}

// kafkaClusterGetCmd represents the get cluster command
var kafkaClusterGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a kafka cluster",
	Long: `Get detail a kafka cluster with kafka cluster ID as input
	Example: bizfly kafka clusters get fd554aac-9ab1-11ea-b09d-bbaf82f02f58`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)

		cluster, err := client.Kafka.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Cluster %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		for _, node := range cluster.Nodes {
			var LanIP []string
			for _, lan := range node.Address.LAN {
				LanIP = append(LanIP, fmt.Sprintf("%s:9093", lan.IPAddress))
			}
			LanIPAddrs := strings.Join(LanIP, ", ")
			var WanIP []string
			for _, wanv4 := range node.Address.WAN_V4 {
				WanIP = append(WanIP, fmt.Sprintf("%s:9094", wanv4.IPAddress))
			}
			WanIPAddrs := strings.Join(WanIP, ", ")
			data = append(data, []string{node.ID, node.Name, node.AvailabilityZone, cluster.Flavor, fmt.Sprintf("%f GB/ %d GB", node.Used, node.VolumeSize), LanIPAddrs, WanIPAddrs, cluster.Status, cluster.CreatedAt})
		}
		formatter.Output(kafkaDetailClusterHeader, data)
	},
}

// kafkaClusterCreateCmd represents the create kafka cluster command
var kafkaClusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a kafka cluster",
	Long:  "Create a new kafka cluster, return a task ID of the processing.\nUse: bizfly kafka clusters create --name <cluster-name> --version <kafka-version-id> --vpc-network-id <vpc-network-id> --nodes <number-of-nodes> --volume-size <volume-size-in-GB> --flavor <flavor-name> [--public-access true|false] [--availability-zone <availability-zone>]",
	Run: func(cmd *cobra.Command, arg []string) {
		scr := gobizfly.KafkaInitClusterRequest{
			ClusterName:      clusterName,
			VersionID:        kafkaVersion,
			VPCNetworkID:     kafkaNetworkInterfaces,
			Nodes:            nodes,
			VolumeSize:       kafkaVolumeSize,
			PublicAccess:     publicAccess,
			Flavor:           kafkaFlavor,
			AvailabilityZone: kafkaAvailabilityZone,
		}
		client, ctx := getApiClient(cmd)
		res, err := client.Kafka.Create(ctx, &scr)
		if err != nil {
			fmt.Printf("Create cluster error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Creating cluster: %#v\n", res)

		fmt.Printf("Creating cluster with task id: %v\n", res.TaskID)
	},
}

// kafkaClusterResizeCmd represents the resize flavor or volume command
var kafkaClusterResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a cluster",
	Long: `
		Resize a cluster.
		Use: bizfly kafka clusters resize <kafka-cluster-id> --type flavor --flavor <flavor name>
		or
		bizfly kafka clusters resize <kafka-cluster-id> --type volume --volume-size <size in GB>
		`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify kafka-cluster-id in the command. Use bizfly kafka resize <kafka-cluster-id> --type <flavor|volume> --flavor or --volume-size")
			os.Exit(1)
		}
		resizeReq := &gobizfly.KafkaResizeClusterRequest{}
		switch kafkaResizeType {
		case "flavor":
			if kafkaFlavor == "" {
				fmt.Println("You need to specify --flavor when resizing type is flavor")
				os.Exit(1)
			}
			resizeReq.Type = "flavor"
			resizeReq.Flavor = kafkaFlavor
		case "volume":
			if kafkaVolumeSize <= 0 {
				fmt.Println("You need to specify --volume-size greater than 0 when resizing type is volume")
				os.Exit(1)
			}
			resizeReq.Type = "volume"
			resizeReq.VolumeSize = kafkaVolumeSize
		default:
			fmt.Println("Invalid type. Use 'flavor' or 'volume'.")
			os.Exit(1)
		}
		kafkaClusterID := args[0]
		client, ctx := getApiClient(cmd)
		res, err := client.Kafka.Resize(ctx, kafkaClusterID, resizeReq)
		if err != nil {
			fmt.Printf("Resize cluster error %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Resizing cluster: %s\n", res.TaskID)
	},
}

// kafkaClusterAddNodeCmd represents the add node command
var kafkaClusterAddNodeCmd = &cobra.Command{
	Use:   "add-node",
	Short: "Add node to a cluster",
	Long: `
		Add node to a cluster.
		Use: bizfly kafka clusters add-node <kafka-cluster-id> --nodes <number of nodes>
		`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to specify kafka-cluster-id in the command. Use bizfly kafka add-node <kafka-cluster-id> --nodes")
			os.Exit(1)
		}
		kafkaClusterID := args[0]
		client, ctx := getApiClient(cmd)
		reqBody := &gobizfly.KafkaAddNodeRequest{
			Nodes: nodes,
			Type:  "increase",
		}
		res, err := client.Kafka.AddNode(ctx, kafkaClusterID, reqBody)
		if err != nil {
			fmt.Printf("Add node error %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Adding node to cluster with task id: %s\n", res.TaskID)
	},
}

// kafkaClusterDeleteCmd represents the delete command
var kafkaClusterDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Kafka Cluster",
	Long: `Delete Kafka Cluster with kafka ID as input
	Example: bizfly kafka clusters delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid arguments")
			_ = cmd.Help() // Display the help message
			return
		}
		client, ctx := getApiClient(cmd)
		_, err := client.Kafka.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Kafka cluster %s is not found", args[0])
				os.Exit(1)
			} else {
				fmt.Printf("Error when get kafka cluster info: %v", err)
				return
			}
		}
		task, err := client.Kafka.Delete(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Kafka cluster %s is not found", args[0])
				os.Exit(1)
			} else {
				fmt.Printf("Error when delete kafka cluster %v", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Deleting kafka cluster with task id: %s\n", task.TaskID)
	},
}

var kafkaFlavorCmd = &cobra.Command{
	Use:   "flavors",
	Short: "Kafka flavors Interaction",
	Long:  `Kafka flavors Interaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // Display the help message
	},
}

var kafkaFlavorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kafka flavors",
	Long:  `List all available Kafka flavors.\nUse: bizfly kafka flavors list`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		flavors, err := client.Kafka.ListFlavor(ctx, nil)
		if err != nil {
			fmt.Printf("Error listing Kafka flavors: %v\n", err)
			os.Exit(1)
		}
		var data [][]string
		for _, flavor := range flavors {
			data = append(data, []string{flavor.ID, flavor.Name, strings.Join([]string{strconv.Itoa(flavor.VCPUs), "Core(s)"}, " "), strings.Join([]string{strconv.Itoa(flavor.RAM), "MB"}, " "), strings.Join([]string{strconv.Itoa(flavor.Disk), "GB"}, " "), flavor.FlavorType})
		}
		formatter.Output(kafkaListFlavorHeader, data)
	},
}

var kafkaVersionCmd = &cobra.Command{
	Use:   "versions",
	Short: "Kafka versions Interaction",
	Long:  `Kafka versions Interaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // Display the help message
	},
}

var kafkaVersionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kafka versions",
	Long:  `List all available Kafka versions.\nUse: bizfly kafka versions list`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		versions, err := client.Kafka.ListVersion(ctx, nil)
		if err != nil {
			fmt.Printf("Error listing Kafka versions: %v\n", err)
			os.Exit(1)
		}
		var data [][]string
		for _, version := range versions {
			data = append(data, []string{version.ID, version.Name, version.Code, strconv.FormatBool(version.IsDefault)})
		}
		formatter.Output(kafkaListVersionHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)
	kafkaCmd.AddCommand(kafkaClusterCmd)
	kafkaClusterCmd.AddCommand(kafkaClusterListCmd)
	kafkaClusterCmd.AddCommand(kafkaClusterGetCmd)

	kafkaCreateFlag := kafkaClusterCreateCmd.Flags()
	kafkaCreateFlag.StringVar(&clusterName, "name", "", "Name of cluster")
	kafkaCreateFlag.StringVar(&kafkaVersion, "version", "", "Version of cluster.\nUse 'bizfly kafka versions list' to see available versions")
	kafkaCreateFlag.StringVar(&kafkaNetworkInterfaces, "vpc-network-id", "", "VPC Network ID.\nUse 'bizfly vpc list' to see available VPC Networks")
	kafkaCreateFlag.IntVar(&kafkaVolumeSize, "volume-size", 0, "Volume Size in GB. Default is 10 GB")
	kafkaCreateFlag.IntVar(&nodes, "nodes", 3, "Number of nodes")
	kafkaCreateFlag.BoolVar(&publicAccess, "public-access", true, "Enable public access")
	kafkaCreateFlag.StringVar(&kafkaFlavor, "flavor", "", "Flavor of cluster.\nUse 'bizfly kafka flavors list' to see available flavors.")
	kafkaCreateFlag.StringVar(&kafkaAvailabilityZone, "availability-zone", "HN1", "Availability zone of cluster")
	_ = kafkaClusterCreateCmd.MarkFlagRequired("name")
	_ = kafkaClusterCreateCmd.MarkFlagRequired("version")
	_ = kafkaClusterCreateCmd.MarkFlagRequired("flavor")
	_ = kafkaClusterCreateCmd.MarkFlagRequired("volume-size")
	_ = kafkaClusterCreateCmd.MarkFlagRequired("vpc-network-id")
	kafkaClusterCmd.AddCommand(kafkaClusterCreateCmd)

	// required one of --flavor or --volume-size
	kafkaClusterResizeCmd.Flags().StringVar(&kafkaResizeType, "type", "", "Type of resizing: flavor or volume")
	_ = kafkaClusterResizeCmd.MarkFlagRequired("type")
	kafkaClusterResizeCmd.Flags().StringVar(&kafkaFlavor, "flavor", "", "New flavor for resizing.\nUse 'bizfly kafka flavors list' to see available flavors")
	kafkaClusterResizeCmd.Flags().IntVar(&kafkaVolumeSize, "volume-size", 0, "New volume size for resizing")
	kafkaClusterCmd.AddCommand(kafkaClusterResizeCmd)

	kafkaClusterCmd.AddCommand(kafkaClusterAddNodeCmd)
	kafkaClusterAddNodeCmd.Flags().StringVar(&kafkaSizeType, "type", "increase", "Type of node addition, e.g. increase|decrease. Default is increase")
	kafkaClusterAddNodeCmd.Flags().IntVar(&nodes, "nodes", 1, "Number of nodes to add.\nDefault is 1")
	_ = kafkaClusterAddNodeCmd.MarkFlagRequired("nodes")
	kafkaClusterCmd.AddCommand(kafkaClusterDeleteCmd)

	kafkaCmd.AddCommand(kafkaFlavorCmd)
	kafkaFlavorCmd.AddCommand(kafkaFlavorListCmd)

	kafkaCmd.AddCommand(kafkaVersionCmd)
	kafkaVersionCmd.AddCommand(kafkaVersionListCmd)
}
