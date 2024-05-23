/*
Copyright © (2020-2022) Bizfly Cloud

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
	lbListHeader            = []string{"ID", "Name", "Network Type", "IP Address", "Operating Status", "Type"}
	poolListHeader          = []string{"ID", "Name", "Algorithm", "Protocol", "Operating Status"}
	listenerListHeader      = []string{"ID", "Name", "Protocol", "Protocol Port", "Operating Status", "Default Pool ID"}
	healthMonitorListHeader = []string{"ID", "Name", "Type", "Delay", "Max Retries", "Timeout", "Operating Status",
		"Domain Name", "URL Path"}
	lbName                          string
	lbType                          string
	networkType                     string
	listenerName                    string
	listenerProtocol                string
	listenerProtocolPort            int
	tlsRef                          string
	lbAlgorithm                     string
	healthMonitorName               string
	healthMonitorDelay              int
	healthMonitorMaxRetries         int
	healthMonitorTimeout            int
	healthMonitorExpectedStatusCode string
	healthMonitorProtocol           string
	healthMonitorURLPath            string
	healthMonitorMaxRetriesDown     int
	healthMonitorMethod             string
	listenerPoolName                string
	defaultPoolID                   string
	poolName                        string
	listenerID                      string
	protocol                        string
	sessionPersistenceType          string
	sessionPersistenceCookieName    string
)

// serverCmd represents the server command
var lbCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "Bizfly Cloud Load Balancer Interaction",
	Long:  `Bizfly Cloud Load Balancer Action: Create, List, Delete`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("loadbalancer called")
	},
}

var lbListenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Bizfly Cloud Load Balancer Listener Interaction",
	Long:  "Bizfly Cloud Load Balancer Listener Action: Create, List, Delete, Get",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var lbPoolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Bizfly Cloud Load Balancer Pool Interaction",
	Long:  "Bizfly Cloud Load Balancer Pool Action: Create, List, Delete, Get",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var lbHealthMonitorCmd = &cobra.Command{
	Use:   "health-monitor",
	Short: "Bizfly Cloud Load Balancer Health Monitor Interaction",
	Long:  "Bizfly Cloud Load Balancer Health Monitor Action: Create, List, Delete, Get",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("health-monitor called")
	},
}

// lbCreateCmd represents the load balancer create command
var lbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a load balancer",
	Long: `Create a load balancer
Example: bizflyctl loadbalancer create --name lb1 --type large --network-type external --listener 8080:8080 --listener 8443:8443 --pool-id pool1 --pool-id pool2 --health-monitor-id hm1`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		payload := gobizfly.LoadBalancerCreateRequest{
			Name:         lbName,
			Type:         lbType,
			VPCNetworkID: vpcNetworkId,
			NetworkType:  networkType,
		}
		if description != "" {
			payload.Description = description
		}
		payload.Listeners = []gobizfly.LoadBalancerListener{}
		payload.Listeners = append(payload.Listeners, gobizfly.LoadBalancerListener{
			Name:          listenerName,
			Protocol:      listenerProtocol,
			DefaultTLSRef: tlsRef,
			ProtocolPort:  listenerProtocolPort,
			DefaultPool: gobizfly.ListenerPool{
				LbAlgorithm: lbAlgorithm,
				Name:        listenerPoolName,
				Protocol:    listenerProtocol,
				Members:     []string{},
				HealthMonitor: gobizfly.ListenerHealthMonitor{
					Delay:          healthMonitorDelay,
					MaxRetries:     healthMonitorMaxRetries,
					Timeout:        healthMonitorTimeout,
					ExpectedCodes:  healthMonitorExpectedStatusCode,
					URLPath:        healthMonitorURLPath,
					MaxRetriesDown: healthMonitorMaxRetriesDown,
					Type:           healthMonitorProtocol,
					HTTPMethod:     healthMonitorMethod,
				},
			},
		})

		lb, err := client.CloudLoadBalancer.Create(ctx, &payload)
		if err != nil {
			log.Fatalf("Error creating load balancer: %v", err)
		}
		var data [][]string
		data = append(data, []string{lb.ID, lb.Name, lb.NetworkType, lb.VipAddress, lb.OperatingStatus, lb.Type})
		formatter.Output(lbListHeader, data)
	},
}

// lbDeleteCmd represents the load balancer delete command
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
			err := client.CloudLoadBalancer.Delete(ctx, &lbdr)
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
		lbs, err := client.CloudLoadBalancer.List(ctx, &gobizfly.ListOptions{})
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

		lb, err := client.CloudLoadBalancer.Get(ctx, args[0])
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
		err := client.CloudLoadBalancer.Pools().Delete(ctx, poolID)
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
		pools, err := client.CloudLoadBalancer.Pools().List(ctx, lbID, &gobizfly.ListOptions{})
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

var lbPoolCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pool in a load balancer",
	Long: `Create a pool in a load balancer
Example: bizfly loadbalancer pool create <loadbalancer_id> --name <pool_name> --protocol <protocol> --lb-algorithm <lb_algorithm>
...`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		payload := &gobizfly.PoolCreateRequest{
			Name:        &poolName,
			Protocol:    protocol,
			LBAlgorithm: lbAlgorithm,
			ListenerID:  listenerID,
		}
		if sessionPersistenceType != "" {
			payload.SessionPersistence = &gobizfly.SessionPersistence{
				Type: sessionPersistenceType,
			}
			if sessionPersistenceType == "APP_COOKIE" {
				payload.SessionPersistence.CookieName = &sessionPersistenceCookieName
			}
		}
		pool, err := client.CloudLoadBalancer.Pools().Create(ctx, args[0], payload)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{pool.ID, pool.Name, pool.LBAlgorithm, pool.Protocol, pool.OperatingStatus})
		formatter.Output(poolListHeader, data)
	},
}

var lbListenerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a listener in a load balancer",
	Long: `Update a listener in a load balancer
Example: bizfly loadbalancer listener update <loadbalancer_id> --name <listener_name> --protocol <protocol> --port <port>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		listener, err := client.CloudLoadBalancer.Listeners().Update(ctx, args[0], &gobizfly.ListenerUpdateRequest{
			Name:                   &listenerName,
			Description:            &description,
			DefaultPoolID:          &defaultPoolID,
			DefaultTLSContainerRef: &tlsRef,
		})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{listener.ID, listener.Name, listener.Protocol, strconv.Itoa(listener.ProtocolPort), listener.OperatingStatus, listener.DefaultPoolID})
		formatter.Output(listenerListHeader, data)
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

		pool, err := client.CloudLoadBalancer.Pools().Get(ctx, args[0])
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

var lbListenerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a listener in a load balancer",
	Long: `Create a listener in a load balancer
Example: bizfly loadbalancer listener create <loadbalancer_id> --name <listener_name> --protocol <protocol> --port <port>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		listener, err := client.CloudLoadBalancer.Listeners().Create(ctx, args[0], &gobizfly.ListenerCreateRequest{
			Name:          &listenerName,
			Description:   &description,
			Protocol:      listenerProtocol,
			ProtocolPort:  listenerProtocolPort,
			DefaultPoolID: &defaultPoolID,
		})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{listener.ID, listener.Name, listener.Protocol, strconv.Itoa(listener.ProtocolPort), listener.OperatingStatus, listener.DefaultPoolID})
		formatter.Output(listenerListHeader, data)
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
		err := client.CloudLoadBalancer.Listeners().Delete(ctx, listenerID)
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
		listeners, err := client.CloudLoadBalancer.Listeners().List(ctx, lbID, &gobizfly.ListOptions{})
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

		listener, err := client.CloudLoadBalancer.Listeners().Get(ctx, args[0])
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

// lbHealthMonitorGetCmd represents the list health monitor command
var lbHealthMonitorGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get health monitor of a listener",
	Long: `Get health monitor of a listener with listener ID as input
Example: bizfly loadbalancer listener get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		healthMontior, err := client.CloudLoadBalancer.HealthMonitors().Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Health monitor of listener %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{healthMontior.ID, healthMontior.Name, healthMontior.Type,
			strconv.Itoa(healthMontior.Delay), strconv.Itoa(healthMontior.TimeOut), strconv.Itoa(healthMontior.MaxRetries),
			healthMontior.DomainName, healthMontior.UrlPath})
		formatter.Output(healthMonitorListHeader, data)
	},
}

var lbHealthMonitorDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete health monitor of a listener",
	Long: `Delete health monitor of a listener with listener ID as input
Example: bizfly loadbalancer listener delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		err := client.CloudLoadBalancer.HealthMonitors().Delete(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Health monitor of listener %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		fmt.Printf("Health monitor of listener %s deleted.", args[0])
	},
}

var lbHealthMonitorCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create health monitor of a listener",
	Long: `Create health monitor of a listener with listener ID as input
Example: bizfly loadbalancer listener create <pool-id> --name sadjf --type HTTP --delay 10 --timeout 10 --max-retries 3 --domain-name www.google.com --url-path /
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		payload := gobizfly.HealthMonitorCreateRequest{
			Name:           healthMonitorName,
			Type:           healthMonitorProtocol,
			Delay:          healthMonitorDelay,
			TimeOut:        healthMonitorTimeout,
			MaxRetries:     healthMonitorMaxRetries,
			URLPath:        healthMonitorURLPath,
			PoolID:         args[0],
			MaxRetriesDown: healthMonitorMaxRetriesDown,
			HTTPMethod:     healthMonitorMethod,
			ExpectedCodes:  healthMonitorExpectedStatusCode,
		}
		healthMonitor, err := client.CloudLoadBalancer.HealthMonitors().Create(ctx, args[0], &payload)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{healthMonitor.ID, healthMonitor.Name, healthMonitor.Type,
			strconv.Itoa(healthMonitor.Delay), strconv.Itoa(healthMonitor.TimeOut), strconv.Itoa(healthMonitor.MaxRetries),
			healthMonitor.DomainName, healthMonitor.UrlPath})
		formatter.Output(healthMonitorListHeader, data)
	},
}

var lbHealthMonitorUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update health monitor of a listener",
	Long: `Update health monitor of a listener with listener ID as input
Example: bizfly loadbalancer listener update <health-monitor-id> --name sadjf --type HTTP --delay 10 --timeout 10 --max-retries 3 --domain-name www.google.com --url-path /`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		payload := gobizfly.HealthMonitorUpdateRequest{
			Name:           healthMonitorName,
			Delay:          &healthMonitorDelay,
			TimeOut:        &healthMonitorTimeout,
			MaxRetries:     &healthMonitorMaxRetries,
			URLPath:        &healthMonitorURLPath,
			MaxRetriesDown: &healthMonitorMaxRetriesDown,
			HTTPMethod:     &healthMonitorMethod,
			ExpectedCodes:  &healthMonitorExpectedStatusCode,
		}
		healthMonitor, err := client.CloudLoadBalancer.HealthMonitors().Update(ctx, args[0], &payload)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{healthMonitor.ID, healthMonitor.Name, healthMonitor.Type,
			strconv.Itoa(healthMonitor.Delay), strconv.Itoa(healthMonitor.TimeOut), strconv.Itoa(healthMonitor.MaxRetries),
			healthMonitor.DomainName, healthMonitor.UrlPath})
		formatter.Output(healthMonitorListHeader, data)
	},
}

var lbResizeLoadBalancerCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a load balancer",
	Long: `Resize a load balancer with load balancer ID, new type (small, medium, large, xtralarge)  as input
	Example: bizfly loadbalancer resize fd554aac-9ab1-11ea-b09d-bbaf82f02f58 medium`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			fmt.Printf("Unknow variable %s", strings.Join(args[2:], ""))
		}
		client, ctx := getApiClient(cmd)
		lbID := args[0]
		newType := args[1]
		err := client.CloudLoadBalancer.Resize(ctx, lbID, newType)
		if err != nil {
			log.Fatal(err)
		}
		lb, err := client.CloudLoadBalancer.Get(ctx, lbID)
		if err != nil {
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
	lbCmd.AddCommand(lbResizeLoadBalancerCmd)
	lbCmd.AddCommand(lbCreateCmd)
	lcpf := lbCreateCmd.PersistentFlags()
	lcpf.StringVar(&lbName, "name", "", "Name of the load balancer")
	lcpf.StringVar(&description, "description", "", "Description of the load balancer")
	lcpf.StringVar(&lbType, "type", "medium", "Type of the load balancer (small, medium, large)")
	_ = cobra.MarkFlagRequired(lcpf, "name")
	lcpf.StringVar(&networkType, "network-type", "external", "Type of the network (external, internal)")
	lcpf.StringVar(&vpcNetworkId, "network-id", "", "ID of the network")
	lcpf.StringVar(&listenerName, "listener-name", "Default Listener", "Name of the listener")
	lcpf.StringVar(&listenerPoolName, "pool-name", "Default", "Name of the pool")
	lcpf.StringVar(&listenerProtocol, "listener-protocol", "HTTP", "Protocol of the listener")
	lcpf.IntVar(&listenerProtocolPort, "listener-port", 80, "Protocol port of the listener")
	lcpf.StringVar(&tlsRef, "tls-ref", "", "TLS reference of the listener")
	lcpf.StringVar(&lbAlgorithm, "algorithm", "ROUND_ROBIN", "Algorithm of the load balancer (ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP)")
	lcpf.IntVar(&healthMonitorDelay, "delay", 5, "Delay of the health monitor")
	lcpf.IntVar(&healthMonitorTimeout, "timeout", 5, "Timeout of the health monitor")
	lcpf.IntVar(&healthMonitorMaxRetries, "max-retries", 3, "Max retries of the health monitor")
	lcpf.StringVar(&healthMonitorExpectedStatusCode, "expected-status-code", "200", "Expected status code of the health monitor")
	lcpf.StringVar(&healthMonitorProtocol, "health-monitor-protocol", "HTTP", "Protocol of the health monitor")
	lcpf.StringVar(&healthMonitorURLPath, "health-monitor-url-path", "/", "URL path of the health monitor")
	lcpf.IntVar(&healthMonitorMaxRetriesDown, "max-retries-down", 3, "Max retries down of the health monitor")
	lcpf.StringVar(&healthMonitorMethod, "health-monitor-method", "GET", "Method of the health monitor")

	lbCmd.AddCommand(lbPoolCmd)
	lbPoolCmd.AddCommand(lbPoolGetCmd)
	lbPoolCmd.AddCommand(lbPoolListCmd)
	lbPoolCmd.AddCommand(lbPoolDeleteCmd)
	lbPoolCmd.AddCommand(lbPoolCreateCmd)
	pcpf := lbPoolCreateCmd.PersistentFlags()
	pcpf.StringVar(&poolName, "name", "", "Name of the load balancer pool")
	_ = cobra.MarkFlagRequired(pcpf, "name")
	pcpf.StringVar(&protocol, "protocol", "HTTP", "Protocol of the pool")
	_ = cobra.MarkFlagRequired(pcpf, "protocol")
	pcpf.StringVar(&lbAlgorithm, "algorithm", "ROUND_ROBIN", "Algorithm of the pool (ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP)")
	pcpf.StringVar(&listenerID, "listener-id", "", "ID of the listener")
	pcpf.StringVar(&sessionPersistenceType, "session-persistence-type", "", "Type of the session persistence (HTTP_COOKIE, APP_COOKIE)")
	pcpf.StringVar(&sessionPersistenceCookieName, "cookie-name", "", "Name of the cookie")

	lbCmd.AddCommand(lbListenerCmd)
	lbListenerCmd.AddCommand(lbListenerGetCmd)
	lbListenerCmd.AddCommand(lbListenerDeleteCmd)
	lbListenerCmd.AddCommand(lbListenerListCmd)
	lbListenerCmd.AddCommand(lbListenerCreateCmd)
	lcf := lbListenerCreateCmd.PersistentFlags()
	lcf.StringVar(&listenerName, "name", "", "Name of the load balancer listener")
	_ = cobra.MarkFlagRequired(lcf, "name")
	lcf.StringVar(&listenerProtocol, "protocol", "HTTP", "Protocol of the listener")
	lcf.IntVar(&listenerProtocolPort, "protocol-port", 80, "Protocol port of the listener")
	lcf.StringVar(&defaultPoolID, "default-pool-id", "", "ID of the default pool")
	lcf.StringVar(&description, "description", "", "Description of the listener")
	lbListenerCmd.AddCommand(lbListenerUpdateCmd)
	lpuf := lbListenerUpdateCmd.PersistentFlags()
	lpuf.StringVar(&listenerName, "name", "", "Name of the load balancer listener")
	lpuf.StringVar(&description, "description", "", "Description of the listener")
	lpuf.StringVar(&defaultPoolID, "default-pool-id", "", "ID of the default pool")
	lpuf.StringVar(&tlsRef, "tls-ref", "", "TLS reference of the listener")

	lbCmd.AddCommand(lbHealthMonitorCmd)
	lbHealthMonitorCmd.AddCommand(lbHealthMonitorGetCmd)
	lbHealthMonitorCmd.AddCommand(lbHealthMonitorDeleteCmd)
	lbHealthMonitorCmd.AddCommand(lbHealthMonitorCreateCmd)
	lhcf := lbHealthMonitorCreateCmd.PersistentFlags()
	lhcf.StringVar(&healthMonitorName, "name", "", "Name of the health monitor")
	_ = cobra.MarkFlagRequired(lhcf, "name")
	lhcf.StringVar(&healthMonitorProtocol, "type", "HTTP", "Type of the health monitor (HTTP, HTTPS, UDP, TCP)")
	lhcf.IntVar(&healthMonitorDelay, "delay", 5, "Delay of the health monitor")
	lhcf.IntVar(&healthMonitorTimeout, "timeout", 5, "Timeout of the health monitor")
	lhcf.IntVar(&healthMonitorMaxRetries, "max-retries", 3, "Max retries of the health monitor")
	lhcf.StringVar(&healthMonitorURLPath, "url-path", "/", "URL path of the health monitor")
	lhcf.StringVar(&healthMonitorExpectedStatusCode, "expected-codes", "", "Expected codes of the health monitor")
	lhcf.StringVar(&healthMonitorMethod, "method", "GET", "Method of the health monitor")
	lhcf.IntVar(&healthMonitorMaxRetriesDown, "max-retries-down", 3, "Max retries down of the health monitor")
	lbHealthMonitorCmd.AddCommand(lbHealthMonitorUpdateCmd)
	lhuf := lbHealthMonitorUpdateCmd.PersistentFlags()
	lhuf.StringVar(&healthMonitorName, "name", "", "Name of the health monitor")
	lhuf.StringVar(&healthMonitorProtocol, "type", "HTTP", "Type of the health monitor (HTTP, HTTPS, UDP, TCP)")
	lhuf.IntVar(&healthMonitorDelay, "delay", 5, "Delay of the health monitor")
	lhuf.IntVar(&healthMonitorTimeout, "timeout", 5, "Timeout of the health monitor")
	lhuf.IntVar(&healthMonitorMaxRetries, "max-retries", 3, "Max retries of the health monitor")
	lhuf.StringVar(&healthMonitorURLPath, "url-path", "/", "URL path of the health monitor")
	lhuf.StringVar(&healthMonitorExpectedStatusCode, "expected-codes", "", "Expected codes of the health monitor")
	lhuf.StringVar(&healthMonitorMethod, "method", "GET", "Method of the health monitor")
	lhuf.IntVar(&healthMonitorMaxRetriesDown, "max-retries-down", 3, "Max retries down of the health monitor")
}
