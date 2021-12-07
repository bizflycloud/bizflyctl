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
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var (
	receiverMethodSupportVerify = []string{
		"sms",
		"telegram",
		"email",
	}
	alarmLoadBalancersTarget = []string{
		"frontend",
		"backend",
	}
	// Use generate tables
	agentListHeader    = table.Row{"ID", "Name", "Runtime", "Online"}
	alarmListHeader    = table.Row{"ID", "Name", "Alarm Type", "Enable"}
	historyListHeader  = table.Row{"ID", "Resource", "Alarm ID", "Alarm Type", "State", "When"}
	receiverListHeader = table.Row{"ID", "Name"}
	resourceGetHeader  = table.Row{"Field", "Value"}
	secretListHeader   = table.Row{"ID", "Name"}

	// Command alarm flags
	alarmAlertInterval    int
	alarmAutoScaling      string
	alarmClusterID        string
	alarmClusterName      string
	alarmComparison       string
	alarmHostname         string
	alarmHTTPExpectedCode int
	alarmHTTPURL          string
	alarmInstances        []string
	alarmLoadBalancers    []string
	alarmName             string
	alarmReceivers        []string
	alarmResourceType     string
	alarmVolumes          []string

	// Command receiver flags
	receiverAutoScaling    string
	receiverEmail          string
	receiverMethods        []string
	receiverName           string
	receiverPhoneNumber    string
	receiverSlack          []string
	receiverTelegramChatID string
	receiverType           string
	receiverWebhook        string

	// Command secret flags
	secretName string
)

func init() {
	rootCmd.AddCommand(alertServiceCmd)

	// Agents
	alertServiceCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentDeleteCmd)
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentShowCmd)

	// Alarms
	alertServiceCmd.AddCommand(alarmCmd)
	alarmCmd.AddCommand(alarmDeleteCmd)
	alarmCmd.AddCommand(alarmDisableCmd)
	alarmCmd.AddCommand(alarmEnableCmd)
	alarmCmd.AddCommand(alarmListCmd)
	alarmCmd.AddCommand(alarmShowCmd)

	// Alarms create
	alarmCreateFlags := alarmCreateCmd.Flags()
	alarmCreateFlags.IntVar(&alarmAlertInterval, "alert-interval", 300, "Time to do sent alert again")
	alarmCreateFlags.IntVar(&alarmHTTPExpectedCode, "http-expected-code", 200, "Expected response status code from API need monitor")
	alarmCreateFlags.StringArrayVarP(&alarmInstances, "instances", "i", []string{}, "Instances need monitor. Example:\n --instances instance-id --instances instance-id-2\nfor multiple instances")
	alarmCreateFlags.StringArrayVarP(&alarmLoadBalancers, "loadbalancers", "l", []string{}, "Load Balancers need monitor. Example:\n --loadbalancers id=fake-id&tgid=fake-target-id&tgtype=frontend\ntgtype maybe:\n - frontend\n - backend")
	alarmCreateFlags.StringArrayVarP(&alarmReceivers, "receivers", "r", []string{}, "Receivers use to received alert. Example:\n --receivers \"id=receiver-id&methods=email,telegram\" --receivers \"id=receiver-id2&methods=method1,method2\"\nfor multiple receivers")
	alarmCreateFlags.StringArrayVarP(&alarmVolumes, "volumes", "v", []string{}, "Volumes need monitor. Example:\n --volumes volume-id --volumes volume-id-2\nfor multiple volumes")
	alarmCreateFlags.StringVar(&alarmClusterID, "cluster-id", "", "Name of cluster use to monitor")
	alarmCreateFlags.StringVar(&alarmClusterName, "cluster-name", "", "ID of cluster use to monitor")
	alarmCreateFlags.StringVar(&alarmHostname, "hostname", "", "Time to do sent alert again")
	alarmCreateFlags.StringVar(&alarmHTTPURL, "http-url", "", "HTTP/HTTPS API need monitor")
	alarmCreateFlags.StringVar(&alarmResourceType, "resource-type", "instance", "Type of of resource to alarm. Maybe include:\n - autoscale_group\n - host\n - http\n - instance\n - load_balancer\n - simple_storage\n - volume")
	alarmCreateFlags.StringVarP(&alarmAutoScaling, "autoscaling", "a", "", "Auto Scaling need monitor. Example:\n --autoscaling id=fake-id&name=fake-name")
	alarmCreateFlags.StringVarP(&alarmComparison, "comparison", "c", "", "Comparison of alarm. Example:\n --comparison \"{\"measurement\":\"iops\",\"compare_type\":\">=\",\"value\":200,\"range_time\":300}\"\n - 'measurement' value maybe:\n    - cpu_used\n    - net_used\n    - ram_used\n  for 'resource_type' is 'instance', 'autoscale_group'.\n    - disk_used\n    - disk_used_percent\n    - iops\n    - read_bytes\n    - write_bytes\n  for 'resource_type' is 'volume'.\n    - request_per_second\n    - data_transfer\n  for 'resource_type' is 'load_balancer'.")
	alarmCreateFlags.StringVarP(&alarmName, "name", "n", "", "Name of alarm")
	err := alarmCreateCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatal(err)
	}
	err = alarmCreateCmd.MarkFlagRequired("receivers")
	if err != nil {
		log.Fatal(err)
	}
	alarmCmd.AddCommand(alarmCreateCmd)

	// Alarms set
	alarmUpdateFlags := alarmSetCmd.Flags()
	alarmUpdateFlags.IntVar(&alarmAlertInterval, "alert-interval", 300, "Time to do sent alert again")
	alarmUpdateFlags.IntVar(&alarmHTTPExpectedCode, "http-expected-code", 200, "Expected response status code from API need monitor")
	alarmUpdateFlags.StringArrayVarP(&alarmInstances, "instances", "i", []string{}, "Instances need monitor. Example:\n --instances instance-id --instances instance-id-2\nfor multiple instances")
	alarmUpdateFlags.StringArrayVarP(&alarmLoadBalancers, "loadbalancers", "l", []string{}, "Load Balancers need monitor. Example:\n --loadbalancers id=fake-id&tgid=fake-target-id&tgtype=frontend\ntgtype maybe:\n - frontend\n - backend")
	alarmUpdateFlags.StringArrayVarP(&alarmReceivers, "receivers", "r", []string{}, "Receivers use to received alert. Example:\n --receivers \"id=receiver-id&methods=email,telegram\" --receivers \"id=receiver-id2&methods=method1,method2\"\nfor multiple receivers")
	alarmUpdateFlags.StringArrayVarP(&alarmVolumes, "volumes", "v", []string{}, "Volumes need monitor. Example:\n --volumes volume-id --volumes volume-id-2\nfor multiple volumes")
	alarmUpdateFlags.StringVar(&alarmClusterID, "cluster-id", "", "Name of cluster use to monitor")
	alarmUpdateFlags.StringVar(&alarmClusterName, "cluster-name", "", "ID of cluster use to monitor")
	alarmUpdateFlags.StringVar(&alarmHostname, "hostname", "", "Time to do sent alert again")
	alarmUpdateFlags.StringVar(&alarmHTTPURL, "http-url", "", "HTTP/HTTPS API need monitor")
	alarmUpdateFlags.StringVar(&alarmResourceType, "resource-type", "", "Type of of resource to alarm. Maybe include:\n - autoscale_group\n - host\n - http\n - instance\n - load_balancer\n - simple_storage\n - volume")
	alarmUpdateFlags.StringVarP(&alarmAutoScaling, "autoscaling", "a", "", "Auto Scaling need monitor. Example:\n --autoscaling id=fake-id&name=fake-name")
	alarmUpdateFlags.StringVarP(&alarmComparison, "comparison", "c", "", "Comparison of alarm. Example:\n --comparison \"{\"measurement\":\"iops\",\"compare_type\":\">=\",\"value\":200,\"range_time\":300}\"\n - 'measurement' value maybe:\n    - cpu_used\n    - net_used\n    - ram_used\n  for 'resource_type' is 'instance', 'autoscale_group'.\n    - disk_used\n    - disk_used_percent\n    - iops\n    - read_bytes\n    - write_bytes\n  for 'resource_type' is 'volume'.\n    - request_per_second\n    - data_transfer\n  for 'resource_type' is 'load_balancer'.")
	alarmUpdateFlags.StringVarP(&alarmName, "name", "n", "", "Name of alarm")
	alarmCmd.AddCommand(alarmSetCmd)

	// Receivers
	alertServiceCmd.AddCommand(receiverCmd)
	receiverCmd.AddCommand(receiverDeleteCmd)
	receiverCmd.AddCommand(receiverListCmd)
	receiverCmd.AddCommand(receiverShowCmd)

	receiverGetVerificationLink := receiverGetVerificationLinkCmd.Flags()
	receiverGetVerificationLink.StringVarP(&receiverType, "type", "t", "", "Type of method being received link verification.\nMaybe include: email, sms, telegram")
	err = receiverGetVerificationLinkCmd.MarkFlagRequired("type")
	if err != nil {
		log.Fatal(err)
	}
	receiverCmd.AddCommand(receiverGetVerificationLinkCmd)

	// Receivers create
	receiverCreateFlags := receiverCreateCmd.Flags()
	receiverCreateFlags.StringArrayVarP(&receiverSlack, "slack", "s", []string{}, "(UNSUPPORTED) Specify Slack.")
	receiverCreateFlags.StringVarP(&receiverAutoScaling, "autoscaling", "a", "", "Specify Auto Scaling Webhook")
	receiverCreateFlags.StringVarP(&receiverEmail, "emailaddress", "e", "", "Specify an email address")
	receiverCreateFlags.StringVarP(&receiverName, "name", "n", "", "Specify name of receiver")
	receiverCreateFlags.StringVarP(&receiverPhoneNumber, "phone", "p", "", "Specify a phone number")
	receiverCreateFlags.StringVarP(&receiverTelegramChatID, "telegram", "t", "", "Specify a telegram chat id")
	receiverCreateFlags.StringVarP(&receiverWebhook, "webhook", "w", "", "Specify a webhook to do trigger")
	err = receiverCreateCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatal(err)
	}
	receiverCmd.AddCommand(receiverCreateCmd)

	// Receivers set
	receiverSetFlags := receiverSetCmd.Flags()
	receiverSetFlags.StringArrayVarP(&receiverSlack, "slack", "s", []string{}, "(UNSUPPORTED) Specify Slack.")
	receiverSetFlags.StringVarP(&receiverAutoScaling, "autoscaling", "a", "", "Specify Auto Scaling Webhook")
	receiverSetFlags.StringVarP(&receiverEmail, "emailaddress", "e", "", "Specify an email address")
	receiverSetFlags.StringVarP(&receiverName, "name", "n", "", "Specify name of receiver")
	receiverSetFlags.StringVarP(&receiverPhoneNumber, "phone", "p", "", "Specify a phone number")
	receiverSetFlags.StringVarP(&receiverTelegramChatID, "telegram", "t", "", "Specify a telegram chat id")
	receiverSetFlags.StringVarP(&receiverWebhook, "webhook", "w", "", "Specify a webhook to do trigger")
	receiverCmd.AddCommand(receiverSetCmd)

	// Receivers unset
	receiverUnSetFlags := receiverUnSetCmd.Flags()
	receiverUnSetFlags.StringArrayVarP(&receiverMethods, "methods", "a", []string{}, "Specify a method to remove. Method maybe include: \n - slack \n - autoscaling \n - emailaddress \n - phone \n - telegram \n - webhook")
	err = receiverUnSetCmd.MarkFlagRequired("methods")
	if err != nil {
		log.Fatal(err)
	}
	receiverCmd.AddCommand(receiverUnSetCmd)

	// Histories
	alertServiceCmd.AddCommand(historyCmd)
	historyCmd.AddCommand(historyListCmd)

	// Secret
	alertServiceCmd.AddCommand(secretCmd)
	secretCmd.AddCommand(secretDeleteCmd)
	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretShowCmd)

	// Secret create
	secretCreateFlags := secretCreateCmd.Flags()
	secretCreateFlags.StringVarP(&secretName, "name", "n", "", "Specify name of secret")
	err = secretCreateCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatal(err)
	}
	secretCmd.AddCommand(secretCreateCmd)
}

var alertServiceCmd = &cobra.Command{
	Use:   "cloudwatcher",
	Short: "BizFly Cloud Watcher Interaction",
	Long:  `Interact with Cloud Watcher Service. Allow do CRUD alarms, receivers, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interacting with cloud watcher service")
	},
}

// ===========================================================================
// Start block interaction for agents
// ===========================================================================
// Sub-command agent
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "BizFly Cloud Watcher Interaction with agent resources",
	Long:  `Interact with Cloud Watcher Service. Allow do CRUD alarms, agents, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interacting with cloud watcher service")
	},
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List agents",
	Long:  "List agents in your account",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		agents, err := client.CloudWatcher.Agents().List(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		for _, agent := range agents {
			s := table.Row{agent.ID, agent.Name, agent.Runtime, agent.Online}
			data = append(data, s)
		}
		formatter.SimpleOutput(agentListHeader, data)
	},
}

var agentShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show detail agent",
	Long:  "Show detail agent by agent ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		agent, err := client.CloudWatcher.Agents().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var jsonSecretData = make(map[string]interface{})
		byteData, err := json.Marshal(agent)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonSecretData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonSecretData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var agentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a agent",
	Long:  "Delete a agent by agent ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		err := client.CloudWatcher.Agents().Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Doing delete agent with ID: %v", args[0])
	},
}

// ===========================================================================
// Start block interaction for alarms
// ===========================================================================
// Sub-command alarm
var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "BizFly Cloud Watcher Interaction with alarm resources",
	Long:  `Interact with Cloud Watcher Service. Allow do CRUD alarms, receivers, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interacting with cloud watcher service")
	},
}

var alarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List alarms",
	Long:  "List alarms in your account",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		alarms, err := client.CloudWatcher.Alarms().List(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		for _, alarm := range alarms {
			s := table.Row{alarm.ID, alarm.Name, alarm.ResourceType, fmt.Sprintf("%v", alarm.Enable)}
			data = append(data, s)
		}
		formatter.SimpleOutput(alarmListHeader, data)
	},
}

var alarmCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an alarm",
	Long:  "Create an alarm",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse receivers from raw input
		var rawReceivers = []map[string]interface{}{}
		for _, alarmReceiver := range alarmReceivers {
			var rawReceiver = make(map[string]interface{})
			ss := strings.Split(strings.ReplaceAll(alarmReceiver, " ", ""), "&")
			for _, parentValue := range ss {
				z := strings.Split(parentValue, "=")
				// Validate data from input
				if z[0] == "" {
					log.Fatal("Not found keyword for: ", z[1])
				}
				if len(z[1]) == 0 {
					log.Fatal("Have error value for: ", z[0])
				}

				if strings.Contains(z[1], ",") {
					rawReceiver[z[0]] = strings.Split(z[1], ",")
				} else {
					rawReceiver[z[0]] = z[1]
				}
			}
			if _, ok := rawReceiver["id"]; !ok {
				log.Fatal("id of receiver is required")
			}
			if _, ok := rawReceiver["methods"]; !ok {
				log.Fatal("methods of receiver is required")
			}
			rawReceivers = append(rawReceivers, rawReceiver)

		}

		// Do make []gobizfly.AlarmReceiversUse
		var alarmCreateReceivers = []gobizfly.AlarmReceiversUse{}
		client, ctx := getApiClient(cmd)
		for _, rawReceiver := range rawReceivers {
			receiver, err := client.CloudWatcher.Receivers().Get(ctx, rawReceiver["id"].(string))
			if err != nil {
				log.Fatal(err)
			}

			var acr = gobizfly.AlarmReceiversUse{
				Name:       receiver.Name,
				ReceiverID: receiver.ReceiverID,
			}
			if _, ok := SliceContains(rawReceiver["methods"], "telegram"); ok {
				MethodsReceiverIsNull(receiver.ReceiverID, "telegram", receiver.TelegramChatID)
				acr.TelegramChatID = receiver.TelegramChatID
			}
			if _, ok := SliceContains(rawReceiver["methods"], "email"); ok {
				MethodsReceiverIsNull(receiver.ReceiverID, "email", receiver.EmailAddress)
				acr.EmailAddress = receiver.EmailAddress
			}
			if _, ok := SliceContains(rawReceiver["methods"], "webhook_url"); ok {
				MethodsReceiverIsNull(receiver.ReceiverID, "webhook_url", receiver.WebhookURL)
				acr.WebhookURL = receiver.WebhookURL
			}
			if _, ok := SliceContains(rawReceiver["methods"], "slack"); ok {
				MethodsReceiverIsNull(receiver.ReceiverID, "slack", receiver.Slack.SlackChannelName)
				acr.SlackChannelName = receiver.Slack.SlackChannelName
			}
			if _, ok := SliceContains(rawReceiver["methods"], "sms"); ok {
				MethodsReceiverIsNull(receiver.ReceiverID, "sms", receiver.SMSNumber)
				acr.SMSNumber = receiver.SMSNumber
				elem, ok := rawReceiver["sms_interval"]
				if ok {
					acr.SMSInterval = elem.(int)
				} else {
					acr.SMSInterval = 1
				}
			}
			if _, ok := SliceContains(rawReceiver["methods"], "autoscaling"); ok {
				acr.AutoscaleClusterName = receiver.AutoScale.ClusterName
			}

			alarmCreateReceivers = append(alarmCreateReceivers, acr)
		}

		if len(alarmLoadBalancers) > 1 {
			log.Fatal("UNSUPPORTED multiple load balancers")
		}
		var rawLoadBalancers = []map[string]interface{}{}
		for _, alarmLoadBalancer := range alarmLoadBalancers {
			var rawLoadBalancer = make(map[string]interface{})
			ss := strings.Split(strings.ReplaceAll(alarmLoadBalancer, " ", ""), "&")
			for _, parentValue := range ss {
				z := strings.Split(parentValue, "=")
				// Validate data from input
				if z[0] == "" {
					log.Fatal("Not found keyword for: ", z[1])
				}
				if len(z[1]) == 0 {
					log.Fatal("Have error value for: ", z[0])
				}

				if strings.Contains(z[1], ",") {
					rawLoadBalancer[z[0]] = strings.Split(z[1], ",")
				} else {
					rawLoadBalancer[z[0]] = z[1]
				}
			}
			if _, ok := rawLoadBalancer["id"]; !ok {
				log.Fatal("id of load balancer is required")
			}
			if _, ok := rawLoadBalancer["tgid"]; !ok {
				log.Fatal("id of backend/frontend of load balancer is required")
			}
			if _, ok := rawLoadBalancer["tgtype"]; !ok {
				log.Fatal("type of tgid is required")
			}
			if _, ok := SliceContains(alarmLoadBalancersTarget, rawLoadBalancer["tgtype"]); !ok {
				log.Fatal("type of tgid is unsupported")
			}
			rawLoadBalancers = append(rawLoadBalancers, rawLoadBalancer)

		}
		var alarmLoadBalancersMonitors = []*gobizfly.AlarmLoadBalancersMonitor{}
		for _, rawLoadBalancer := range rawLoadBalancers {
			lb, err := client.LoadBalancer.Get(ctx, rawLoadBalancer["id"].(string))
			if err != nil {
				log.Fatal(err)
			}

			var albm = gobizfly.AlarmLoadBalancersMonitor{
				LoadBalancerID:   lb.ID,
				LoadBalancerName: lb.Name,
				TargetID:         rawLoadBalancer["tgid"].(string),
				TargetType:       rawLoadBalancer["tgtype"].(string),
			}

			// Validate frontend/backend of loadbalancer
			if rawLoadBalancer["tgtype"] == "frontend" {
				frontend, err := client.Listener.Get(ctx, rawLoadBalancer["tgid"].(string))
				if err != nil {
					log.Fatal(err)
				}
				albm.TargetName = frontend.Name
			} else {
				backend, err := client.Pool.Get(ctx, rawLoadBalancer["tgid"].(string))
				if err != nil {
					log.Fatal(err)
				}
				albm.TargetName = backend.Name
			}

			alarmLoadBalancersMonitors = append(alarmLoadBalancersMonitors, &albm)
		}

		var alarmCreateRequest = gobizfly.AlarmCreateRequest{
			AlertInterval:    alarmAlertInterval,
			ClusterID:        alarmClusterID,
			ClusterName:      alarmClusterName,
			Hostname:         alarmHostname,
			HTTPExpectedCode: alarmHTTPExpectedCode,
			HTTPURL:          alarmHTTPURL,
			Name:             alarmName,
			Receivers:        alarmCreateReceivers,
			LoadBalancers:    alarmLoadBalancersMonitors,
			ResourceType:     alarmResourceType,
		}

		if alarmComparison != "" {
			comparison := make(map[string]interface{})
			err := json.Unmarshal([]byte(alarmComparison), &comparison)
			if err != nil {
				log.Fatal(err)
			}

			rangetime, err := strconv.Atoi(fmt.Sprintf("%v", comparison["range_time"]))
			if err != nil {
				log.Fatal(err)
			}
			alarmCreateRequest.Comparison = &gobizfly.Comparison{
				Measurement: comparison["measurement"].(string),
				RangeTime:   rangetime,
				Value:       comparison["value"].(float64),
				CompareType: comparison["compare_type"].(string),
			}
		} else {
			alarmCreateRequest.Comparison = &gobizfly.Comparison{}
		}

		var volumesMonitor = []gobizfly.AlarmVolumesMonitor{}
		if len(alarmVolumes) > 0 {
			for _, volumeID := range alarmVolumes {
				volume, err := client.Volume.Get(ctx, volumeID)
				if err != nil {
					log.Fatal(err)
				}
				volumesMonitor = append(volumesMonitor, gobizfly.AlarmVolumesMonitor{
					ID:   volume.ID,
					Name: volume.Name,
				})
			}
			alarmCreateRequest.Volumes = &volumesMonitor
		}

		var instancesMonitor = []gobizfly.AlarmInstancesMonitors{}
		if len(alarmInstances) > 0 {
			for _, instanceID := range alarmInstances {
				instance, err := client.Server.Get(ctx, instanceID)
				if err != nil {
					log.Fatal(err)
				}
				instancesMonitor = append(instancesMonitor, gobizfly.AlarmInstancesMonitors{
					ID:   instance.ID,
					Name: instance.Name,
				})
			}
			alarmCreateRequest.Instances = &instancesMonitor
		}

		response, err := client.CloudWatcher.Alarms().Create(ctx, &alarmCreateRequest)
		if err != nil {
			log.Fatal(err)
		}
		alarm, _ := client.CloudWatcher.Alarms().Get(ctx, response.ID)

		var jsonAlarmData = make(map[string]interface{})
		byteData, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonAlarmData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonAlarmData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var alarmShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show detail alarm",
	Long:  "Show detail alarm by alarm ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		alarm, err := client.CloudWatcher.Alarms().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var jsonAlarmData = make(map[string]interface{})
		byteData, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonAlarmData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonAlarmData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var alarmDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an alarm",
	Long:  "Delete an alarm by alarm ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		err := client.CloudWatcher.Alarms().Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Doing delete alarm with ID: %v", args[0])
	},
}

var alarmSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Update an alarm",
	Long:  "Update an alarm",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %v", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		oldAlarm, err := client.CloudWatcher.Alarms().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var alarmCreateReceivers = []gobizfly.AlarmReceiversUse{}
		if len(alarmReceivers) > 0 {
			// Parse receivers from raw input
			var rawReceivers = []map[string]interface{}{}
			for _, alarmReceiver := range alarmReceivers {
				var rawReceiver = make(map[string]interface{})
				ss := strings.Split(strings.ReplaceAll(alarmReceiver, " ", ""), "&")
				for _, parentValue := range ss {
					z := strings.Split(parentValue, "=")
					// Validate data from input
					if z[0] == "" {
						log.Fatal("Not found keyword for: ", z[1])
					}
					if len(z[1]) == 0 {
						log.Fatal("Have error value for: ", z[0])
					}

					if strings.Contains(z[1], ",") {
						rawReceiver[z[0]] = strings.Split(z[1], ",")
					} else {
						rawReceiver[z[0]] = z[1]
					}
				}
				if _, ok := rawReceiver["id"]; !ok {
					log.Fatal("id of receiver is required")
				}
				if _, ok := rawReceiver["methods"]; !ok {
					log.Fatal("methods of receiver is required")
				}
				rawReceivers = append(rawReceivers, rawReceiver)

			}

			// Do make []gobizfly.AlarmReceiversUse
			for _, rawReceiver := range rawReceivers {
				receiver, err := client.CloudWatcher.Receivers().Get(ctx, rawReceiver["id"].(string))
				if err != nil {
					log.Fatal(err)
				}

				var acr = gobizfly.AlarmReceiversUse{
					Name:       receiver.Name,
					ReceiverID: receiver.ReceiverID,
				}
				if _, ok := SliceContains(rawReceiver["methods"], "telegram"); ok {
					MethodsReceiverIsNull(receiver.ReceiverID, "telegram", receiver.TelegramChatID)
					acr.TelegramChatID = receiver.TelegramChatID
				}
				if _, ok := SliceContains(rawReceiver["methods"], "email"); ok {
					MethodsReceiverIsNull(receiver.ReceiverID, "email", receiver.EmailAddress)
					acr.EmailAddress = receiver.EmailAddress
				}
				if _, ok := SliceContains(rawReceiver["methods"], "webhook_url"); ok {
					MethodsReceiverIsNull(receiver.ReceiverID, "webhook_url", receiver.WebhookURL)
					acr.WebhookURL = receiver.WebhookURL
				}
				if _, ok := SliceContains(rawReceiver["methods"], "slack"); ok {
					MethodsReceiverIsNull(receiver.ReceiverID, "slack", receiver.Slack.SlackChannelName)
					acr.SlackChannelName = receiver.Slack.SlackChannelName
				}
				if _, ok := SliceContains(rawReceiver["methods"], "sms"); ok {
					MethodsReceiverIsNull(receiver.ReceiverID, "sms", receiver.SMSNumber)
					acr.SMSNumber = receiver.SMSNumber
					elem, ok := rawReceiver["sms_interval"]
					if ok {
						acr.SMSInterval = elem.(int)
					} else {
						acr.SMSInterval = 1
					}
				}
				if _, ok := SliceContains(rawReceiver["methods"], "autoscaling"); ok {
					acr.AutoscaleClusterName = receiver.AutoScale.ClusterName
				}

				alarmCreateReceivers = append(alarmCreateReceivers, acr)
			}
		} else {
			alarmCreateReceivers = oldAlarm.Receivers
		}

		if alarmClusterID == "" {
			alarmClusterID = oldAlarm.ClusterID
		}
		if alarmClusterName == "" {
			alarmClusterName = oldAlarm.ClusterName
		}
		if alarmHostname == "" {
			alarmHostname = oldAlarm.Hostname
		}

		if alarmHTTPURL == "" {
			alarmHTTPURL = oldAlarm.HTTPURL
		}
		if alarmName == "" {
			alarmName = oldAlarm.Name
		}
		if alarmResourceType == "" {
			alarmResourceType = oldAlarm.ResourceType
		}

		var alarmLoadBalancersMonitors = []*gobizfly.AlarmLoadBalancersMonitor{}
		if len(alarmLoadBalancers) > 0 {
			if len(alarmLoadBalancers) > 1 {
				log.Fatal("UNSUPPORTED multiple load balancers")
			}
			var rawLoadBalancers = []map[string]interface{}{}
			for _, alarmLoadBalancer := range alarmLoadBalancers {
				var rawLoadBalancer = make(map[string]interface{})
				ss := strings.Split(strings.ReplaceAll(alarmLoadBalancer, " ", ""), "&")
				for _, parentValue := range ss {
					z := strings.Split(parentValue, "=")
					// Validate data from input
					if z[0] == "" {
						log.Fatal("Not found keyword for: ", z[1])
					}
					if len(z[1]) == 0 {
						log.Fatal("Have error value for: ", z[0])
					}

					if strings.Contains(z[1], ",") {
						rawLoadBalancer[z[0]] = strings.Split(z[1], ",")
					} else {
						rawLoadBalancer[z[0]] = z[1]
					}
				}
				if _, ok := rawLoadBalancer["id"]; !ok {
					log.Fatal("id of load balancer is required")
				}
				if _, ok := rawLoadBalancer["tgid"]; !ok {
					log.Fatal("id of backend/frontend of load balancer is required")
				}
				if _, ok := rawLoadBalancer["tgtype"]; !ok {
					log.Fatal("type of tgid is required")
				}
				if _, ok := SliceContains(alarmLoadBalancersTarget, rawLoadBalancer["tgtype"]); !ok {
					log.Fatal("type of tgid is unsupported")
				}
				rawLoadBalancers = append(rawLoadBalancers, rawLoadBalancer)

			}
			for _, rawLoadBalancer := range rawLoadBalancers {
				lb, err := client.LoadBalancer.Get(ctx, rawLoadBalancer["id"].(string))
				if err != nil {
					log.Fatal(err)
				}

				var albm = gobizfly.AlarmLoadBalancersMonitor{
					LoadBalancerID:   lb.ID,
					LoadBalancerName: lb.Name,
					TargetID:         rawLoadBalancer["tgid"].(string),
					TargetType:       rawLoadBalancer["tgtype"].(string),
				}

				// Validate frontend/backend of loadbalancer
				if rawLoadBalancer["tgtype"] == "frontend" {
					frontend, err := client.Listener.Get(ctx, rawLoadBalancer["tgid"].(string))
					if err != nil {
						log.Fatal(err)
					}
					albm.TargetName = frontend.Name
				} else {
					backend, err := client.Pool.Get(ctx, rawLoadBalancer["tgid"].(string))
					if err != nil {
						log.Fatal(err)
					}
					albm.TargetName = backend.Name
				}

				alarmLoadBalancersMonitors = append(alarmLoadBalancersMonitors, &albm)
			}
		} else {
			alarmLoadBalancersMonitors = oldAlarm.LoadBalancers
		}

		var alarmUpdateRequest = gobizfly.AlarmUpdateRequest{
			AlertInterval:    alarmAlertInterval,
			ClusterID:        alarmClusterID,
			ClusterName:      alarmClusterName,
			Hostname:         alarmHostname,
			HTTPExpectedCode: alarmHTTPExpectedCode,
			HTTPURL:          alarmHTTPURL,
			Name:             alarmName,
			Receivers:        &alarmCreateReceivers,
			ResourceType:     alarmResourceType,
			LoadBalancers:    alarmLoadBalancersMonitors,
		}

		if alarmComparison != "" {
			comparison := make(map[string]interface{})
			err := json.Unmarshal([]byte(alarmComparison), &comparison)
			if err != nil {
				log.Fatal(err)
			}

			rangetime, err := strconv.Atoi(fmt.Sprintf("%v", comparison["range_time"]))
			if err != nil {
				log.Fatal(err)
			}
			alarmUpdateRequest.Comparison = &gobizfly.Comparison{
				Measurement: comparison["measurement"].(string),
				RangeTime:   rangetime,
				Value:       comparison["value"].(float64),
				CompareType: comparison["compare_type"].(string),
			}
		} else {
			alarmUpdateRequest.Comparison = oldAlarm.Comparison
		}

		var volumesMonitor = []gobizfly.AlarmVolumesMonitor{}
		if len(alarmVolumes) > 0 {
			for _, volumeID := range alarmVolumes {
				volume, err := client.Volume.Get(ctx, volumeID)
				if err != nil {
					log.Fatal(err)
				}
				volumesMonitor = append(volumesMonitor, gobizfly.AlarmVolumesMonitor{
					ID:   volume.ID,
					Name: volume.Name,
				})
			}
			alarmUpdateRequest.Volumes = &volumesMonitor
		} else {
			if len(oldAlarm.Volumes) > 0 {
				alarmUpdateRequest.Volumes = &oldAlarm.Volumes
			}
		}

		var instancesMonitor = []gobizfly.AlarmInstancesMonitors{}
		if len(alarmInstances) > 0 {
			for _, instanceID := range alarmInstances {
				instance, err := client.Server.Get(ctx, instanceID)
				if err != nil {
					log.Fatal(err)
				}
				instancesMonitor = append(instancesMonitor, gobizfly.AlarmInstancesMonitors{
					ID:   instance.ID,
					Name: instance.Name,
				})
			}
			alarmUpdateRequest.Instances = &instancesMonitor
		} else {
			if len(oldAlarm.Instances) > 0 {
				alarmUpdateRequest.Instances = &oldAlarm.Instances
			}
		}

		response, err := client.CloudWatcher.Alarms().Update(ctx, args[0], &alarmUpdateRequest)
		if err != nil {
			log.Fatal(err)
		}
		alarm, _ := client.CloudWatcher.Alarms().Get(ctx, response.ID)

		var jsonAlarmData = make(map[string]interface{})
		byteData, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonAlarmData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonAlarmData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var alarmEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable an alarm",
	Long:  "Enable an alarm",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %v", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		alarmUpdateRequest := gobizfly.AlarmUpdateRequest{
			Enable: true,
		}
		response, err := client.CloudWatcher.Alarms().Update(ctx, args[0], &alarmUpdateRequest)
		if err != nil {
			log.Fatal(err)
		}
		alarm, _ := client.CloudWatcher.Alarms().Get(ctx, response.ID)

		var jsonAlarmData = make(map[string]interface{})
		byteData, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonAlarmData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonAlarmData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var alarmDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable an alarm",
	Long:  "Disable an alarm",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %v", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		alarmUpdateRequest := gobizfly.AlarmUpdateRequest{
			Enable: false,
		}
		response, err := client.CloudWatcher.Alarms().Update(ctx, args[0], &alarmUpdateRequest)
		if err != nil {
			log.Fatal(err)
		}
		alarm, _ := client.CloudWatcher.Alarms().Get(ctx, response.ID)

		var jsonAlarmData = make(map[string]interface{})
		byteData, err := json.Marshal(alarm)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonAlarmData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonAlarmData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

// ===========================================================================
// Start block interaction for receivers
// ===========================================================================
// Sub-command receiver
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "BizFly Cloud Watcher Interaction with receiver resources",
	Long:  `Interact with Cloud Watcher Service. Allow do CRUD alarms, receivers, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interacting with cloud watcher service")
	},
}

var receiverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List receivers",
	Long:  "List receivers in your account",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		receivers, err := client.CloudWatcher.Receivers().List(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		for _, receiver := range receivers {
			s := table.Row{receiver.ReceiverID, receiver.Name}
			data = append(data, s)
		}
		formatter.SimpleOutput(receiverListHeader, data)
	},
}

var receiverCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an receiver",
	Long:  "Create an receiver by specific informations",
	Run: func(cmd *cobra.Command, args []string) {
		var rcr = gobizfly.ReceiverCreateRequest{}
		client, ctx := getApiClient(cmd)

		// current not need handle
		// if len(receiverSlack) > 0 {
		// }
		if len(receiverAutoScaling) > 0 {
			// Parse autoscaling from raw input
			var rawAutoScaling = make(map[string]string)
			ss := strings.Split(receiverAutoScaling, "&")
			for _, parentValue := range ss {
				z := strings.Split(parentValue, "=")
				// Validate data from input
				if z[0] == "" {
					log.Fatal("Not found keyword for: ", z[1])
				}
				if len(z[1]) == 0 {
					log.Fatal("Have error value for: ", z[0])
				}

				rawAutoScaling[z[0]] = z[1]
			}
			if _, ok := rawAutoScaling["type"]; !ok {
				log.Fatal("action type is required for auto scaling group")
			}
			if _, ok := rawAutoScaling["id"]; !ok {
				log.Fatal("id of for auto scaling group is required")
			}
			webhook, err := client.AutoScaling.Webhooks().Get(ctx, rawAutoScaling["id"], rawAutoScaling["type"])
			if err != nil {
				log.Fatal(err)
			}
			rcr.AutoScale = webhook
		}

		if len(receiverEmail) > 0 {
			rcr.EmailAddress = receiverEmail
		}
		if len(receiverName) > 0 {
			rcr.Name = receiverName
		}
		if len(receiverPhoneNumber) > 0 {
			rcr.SMSNumber = receiverPhoneNumber
		}
		if len(receiverTelegramChatID) > 0 {
			rcr.TelegramChatID = receiverTelegramChatID
		}
		if len(receiverWebhook) > 0 {
			rcr.WebhookURL = receiverWebhook
		}

		response, err := client.CloudWatcher.Receivers().Create(ctx, &rcr)
		if err != nil {
			log.Fatal(err)
		}

		receiver, err := client.CloudWatcher.Receivers().Get(ctx, response.ID)
		if err != nil {
			log.Fatal(err)
		}

		var jsonReceiverData = make(map[string]interface{})
		byteData, err := json.Marshal(receiver)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonReceiverData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonReceiverData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var receiverGetVerificationLinkCmd = &cobra.Command{
	Use:   "verify",
	Short: "Get a link verify a method of receiver",
	Long:  "Get a link verify a method of receiver by specific informations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %v", strings.Join(args[1:], ""))
		}
		if _, ok := SliceContains(receiverMethodSupportVerify, receiverType); !ok {
			log.Fatalf("Method %v is unsupported to get link verification", receiverType)
		}

		client, ctx := getApiClient(cmd)
		if err := client.CloudWatcher.Receivers().ResendVerificationLink(ctx, args[0], receiverType); err == nil {
			log.Printf("A link verification was sent to %v of receiver %v", receiverType, args[0])
		} else {
			log.Fatalf("Failed to sent link verification to %v of receiver %v", receiverType, args[0])
		}
	},
}

var receiverShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show detail receiver",
	Long:  "Show detail receiver by receiver ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		receiver, err := client.CloudWatcher.Receivers().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var jsonReceiverData = make(map[string]interface{})
		byteData, err := json.Marshal(receiver)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonReceiverData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonReceiverData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var receiverDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a receiver",
	Long:  "Delete a receiver by receiver ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		err := client.CloudWatcher.Receivers().Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Doing delete receiver with ID: %v", args[0])
	},
}

var receiverSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Update an receiver",
	Long:  "Update an receiver by specific informations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %v", strings.Join(args[1:], ""))
		}

		client, ctx := getApiClient(cmd)
		oldReceiver, err := client.CloudWatcher.Receivers().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var rcr = gobizfly.ReceiverCreateRequest{}
		// current not need handle
		// if len(receiverSlack) > 0 {
		// }
		if len(receiverAutoScaling) > 0 {
			// Parse autoscaling from raw input
			var rawAutoScaling = make(map[string]string)
			ss := strings.Split(receiverAutoScaling, "&")
			for _, parentValue := range ss {
				z := strings.Split(parentValue, "=")
				// Validate data from input
				if z[0] == "" {
					log.Fatal("Not found keyword for: ", z[1])
				}
				if len(z[1]) == 0 {
					log.Fatal("Have error value for: ", z[0])
				}

				rawAutoScaling[z[0]] = z[1]
			}
			if _, ok := rawAutoScaling["type"]; !ok {
				log.Fatal("action type is required for auto scaling group")
			}
			if _, ok := rawAutoScaling["id"]; !ok {
				log.Fatal("id of for auto scaling group is required")
			}
			webhook, err := client.AutoScaling.Webhooks().Get(ctx, rawAutoScaling["id"], rawAutoScaling["type"])
			if err != nil {
				log.Fatal(err)
			}
			rcr.AutoScale = webhook
		} else {
			rcr.AutoScale = oldReceiver.AutoScale
		}

		if len(receiverEmail) > 0 {
			rcr.EmailAddress = receiverEmail
		} else {
			rcr.EmailAddress = oldReceiver.EmailAddress
		}
		if len(receiverName) > 0 {
			rcr.Name = receiverName
		} else {
			rcr.Name = oldReceiver.Name
		}
		if len(receiverPhoneNumber) > 0 {
			rcr.SMSNumber = receiverPhoneNumber
		} else {
			rcr.SMSNumber = oldReceiver.SMSNumber
		}
		if len(receiverTelegramChatID) > 0 {
			rcr.TelegramChatID = receiverTelegramChatID
		} else {
			rcr.TelegramChatID = oldReceiver.TelegramChatID
		}
		if len(receiverWebhook) > 0 {
			rcr.WebhookURL = receiverWebhook
		} else {
			rcr.WebhookURL = oldReceiver.WebhookURL
		}

		response, err := client.CloudWatcher.Receivers().Update(ctx, args[0], &rcr)
		if err != nil {
			log.Fatal(err)
		}

		receiver, err := client.CloudWatcher.Receivers().Get(ctx, response.ID)
		if err != nil {
			log.Fatal(err)
		}

		var jsonReceiverData = make(map[string]interface{})
		byteData, err := json.Marshal(receiver)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonReceiverData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonReceiverData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var receiverUnSetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove a method receiver",
	Long:  "Remove a method receiver by specific informations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %v", strings.Join(args[1:], ""))
		}

		client, ctx := getApiClient(cmd)
		oldReceiver, err := client.CloudWatcher.Receivers().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var rcr = gobizfly.ReceiverCreateRequest{
			AutoScale:      oldReceiver.AutoScale,
			EmailAddress:   oldReceiver.EmailAddress,
			Name:           oldReceiver.Name,
			SMSNumber:      oldReceiver.SMSNumber,
			TelegramChatID: oldReceiver.TelegramChatID,
			WebhookURL:     oldReceiver.WebhookURL,
		}
		for _, method := range receiverMethods {
			// current not need handle
			// if method == "slack" {
			// }

			if method == "autoscaling" {
				rcr.AutoScale = nil
			}

			if method == "emailaddress" {
				rcr.EmailAddress = ""
			}

			if method == "phone" {
				rcr.SMSNumber = ""
			}
			if method == "telegram" {
				rcr.TelegramChatID = ""
			}
			if method == "webhook" {
				rcr.WebhookURL = ""
			}
		}

		response, err := client.CloudWatcher.Receivers().Update(ctx, args[0], &rcr)
		if err != nil {
			log.Fatal(err)
		}

		receiver, err := client.CloudWatcher.Receivers().Get(ctx, response.ID)
		if err != nil {
			log.Fatal(err)
		}

		var jsonReceiverData = make(map[string]interface{})
		byteData, err := json.Marshal(receiver)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonReceiverData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonReceiverData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

// ===========================================================================
// Start block interaction for histories
// ===========================================================================
// Sub-command history
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "BizFly Cloud Watcher Interaction with history resources",
	Long:  `Interact with Cloud Watcher Service. Allow do CRUD alarms, receivers, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interacting with cloud watcher service")
	},
}

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List history",
	Long:  "List 26 latest history in your account",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		histories, err := client.CloudWatcher.Histories().List(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		for _, history := range histories {
			s := table.Row{
				history.HistoryID,
				history.Resource.(string),
				history.AlarmID,
				history.Alarm.ResourceType,
				history.State,
				history.Created,
			}
			data = append(data, s)
		}
		formatter.SimpleOutput(historyListHeader, data)
	},
}

// ===========================================================================
// Start block interaction for secrets
// ===========================================================================
// Sub-command secret
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "BizFly Cloud Watcher Interaction with secret resources",
	Long:  `Interact with Cloud Watcher Service. Allow do CRUD alarms, secrets, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interacting with cloud watcher service")
	},
}

var secretListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Long:  "List secrets in your account",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		secrets, err := client.CloudWatcher.Secrets().List(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		for _, secret := range secrets {
			s := table.Row{secret.ID, secret.Name}
			data = append(data, s)
		}
		formatter.SimpleOutput(secretListHeader, data)
	},
}

var secretCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an secret",
	Long:  "Create an secret by specific informations",
	Run: func(cmd *cobra.Command, args []string) {
		var scr = gobizfly.SecretsCreateRequest{}
		client, ctx := getApiClient(cmd)

		if len(secretName) > 0 {
			scr.Name = secretName
		}

		response, err := client.CloudWatcher.Secrets().Create(ctx, &scr)
		if err != nil {
			log.Fatal(err)
		}

		receiver, err := client.CloudWatcher.Secrets().Get(ctx, response.ID)
		if err != nil {
			log.Fatal(err)
		}

		var jsonSecretData = make(map[string]interface{})
		byteData, err := json.Marshal(receiver)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonSecretData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonSecretData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var secretShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show detail secret",
	Long:  "Show detail secret by secret ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		secret, err := client.CloudWatcher.Secrets().Get(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		var jsonSecretData = make(map[string]interface{})
		byteData, err := json.Marshal(secret)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(byteData, &jsonSecretData)
		if err != nil {
			log.Fatal(err)
		}

		var data []table.Row
		data = ProcessDataTables(data, jsonSecretData)
		formatter.SimpleOutput(resourceGetHeader, data)
	},
}

var secretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret",
	Long:  "Delete a secret by secret ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := getApiClient(cmd)
		err := client.CloudWatcher.Secrets().Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Doing delete secret with ID: %v", args[0])
	},
}
