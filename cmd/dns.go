package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

var (
	zoneID              string
	required            bool
	zoneDescription     string
	recordData          string
	recordName          string
	recordType          string
	httpHealthCheck     string
	tcpHealthCheck      string
	ipv4RoutingPolicy   []string
	ipv6RoutingPolicy   []string
	TTL                 int
	domainData          []string
	NormalTypes         = []string{"A", "AAAA", "TXT"}
	NormalDataHeader    = []string{"Data"}
	MXDataHeader        = []string{"Domain", "Priority"}
	PolicyRoutingHeader = []string{"Region", "IPv4", "IPv6"}
	zonesHeader         = []string{"ID", "Name", "Deleted", "NameServer", "TTL", "Active", "Created At", "Update At"}
	recordSetHeader     = []string{"ID", "Name", "Type", "TTL"}
)

type recordPayload struct {
	Record interface{} `json:"record"`
}

var dnsComnmand = &cobra.Command{
	Use:   "dns",
	Short: "Bizfly Cloud DNS Interaction",
	Long:  "Bizfly Cloud DNS Action: List zones, Create zone, Get zone, Delete zone, Create record, Get record, Delete record",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dns called")
	},
}

var listZonesCommand = &cobra.Command{
	Use:   "list-zones",
	Short: "List all zones",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		resp, err := client.DNS.ListZones(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		zones := resp.Zones
		var data [][]string
		for _, zone := range zones {
			nameserverString := strings.Join(zone.NameServer, "\n")
			data = append(data, []string{zone.ID, zone.Name, strconv.Itoa(zone.Deleted),
				nameserverString, strconv.Itoa(zone.TTL), strconv.FormatBool(zone.Active),
				zone.CreatedAt, zone.UpdatedAt})
		}
		formatter.Output(zonesHeader, data)
	},
}

var getZoneCommand = &cobra.Command{
	Use:   "get-zone",
	Short: "Get a zone",
	Long: `Get a zone
Usage: ./bizfly dns get-zone <zone-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		resp, err := client.DNS.GetZone(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		zone := resp.Zone
		recordSets := resp.RecordsSet

		nameserverString := strings.Join(zone.NameServer, "\n")
		var zoneData [][]string
		zoneData = append(zoneData, []string{zone.ID, zone.Name, strconv.Itoa(zone.Deleted),
			nameserverString, strconv.Itoa(zone.TTL), strconv.FormatBool(zone.Active),
			zone.CreatedAt, zone.UpdatedAt})

		var recordSetData [][]string
		for _, recordSet := range recordSets {
			recordSetData = append(recordSetData, []string{recordSet.ID, recordSet.Name,
				recordSet.Type, recordSet.TTL})
		}
		formatter.Output(zonesHeader, zoneData)
		formatter.Output(recordSetHeader, recordSetData)
	},
}

var createZoneCommand = &cobra.Command{
	Use:   "create-zone",
	Short: "Create DNS Zone",
	Long: `Create DNS Zone
Usage: ./bizfly dns create-zone <zone-name> [flags]`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		payload := &gobizfly.CreateZonePayload{
			Name:        args[0],
			Required:    required,
			Description: description,
		}
		resp, err := client.DNS.CreateZone(ctx, payload)
		if err != nil {
			log.Fatal(err)
		}
		zone := resp.Zone
		recordSets := resp.RecordsSet

		nameserverString := strings.Join(zone.NameServer, "\n")
		var zoneData [][]string
		zoneData = append(zoneData, []string{zone.ID, zone.Name, strconv.Itoa(zone.Deleted),
			nameserverString, strconv.Itoa(zone.TTL), strconv.FormatBool(zone.Active),
			zone.CreatedAt, zone.UpdatedAt})

		var recordSetData [][]string
		for _, recordSet := range recordSets {
			recordSetData = append(recordSetData, []string{recordSet.ID, recordSet.Name,
				recordSet.Type, recordSet.TTL})
		}
		formatter.Output(zonesHeader, zoneData)
		formatter.Output(recordSetHeader, recordSetData)
	},
}

var deleteZoneCommand = &cobra.Command{
	Use:   "delete-zone",
	Short: "Delete zone",
	Long: `Delete zone
Usage: ./bizfly dns delete zone <zone-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		err := client.DNS.DeleteZone(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Deleted Zone ", zoneID)
	},
}

var getRecordCommand = &cobra.Command{
	Use:   "get-record",
	Short: "Get record via ID",
	Long: `Get DNS record in a zone
Usage: ./bizfly dns get-record <record-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		recordSet, err := client.DNS.GetRecord(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var recordSetData [][]string
		recordSetData = append(recordSetData, []string{recordSet.ID, recordSet.Name,
			recordSet.Type, strconv.Itoa(recordSet.TTL)})
		formatter.Output(recordSetHeader, recordSetData)
		outputRecordData(recordSet)
	},
}

var createRecordCommand = &cobra.Command{
	Use:   "create-record",
	Short: "Create new record in a zone",
	Long: `Create DNS record in a zone
General arguments:
  - record-name: Name of the record
  - record-type: Type of the record
  - ttl: Time to live
Type arguments:
  - Type: Normal (A, AAAA, TXT)
    + data: IPv4/v6 addresses or domains depends on its type
    Example: ./bizfly dns create-record --zone-id zone-123 --name test_a_1 --ttl 600 --type A --data "123.123.123.123;7.7.7.7
  - Type: MX
    + domain-data: specify the domains and its priority. Format: --domain-data domain:priority
    Example: ./bizfly dns create-record --zone-id 123-zone --name test_mx_1 --ttl 600 --type MX --domain-data test.com:10 --domain-data test1.com:49
  - Type: GEOIP
    + ipv4-policy: Specify IPv4 routing policy in regions (HN, HCM, SG, USA). Format: --ipv4-policy "region_name:<addr1:str>,<addr2:str>" 
    + ipv6-policy: Specify IPv6 routing policy in regions (HN, HCM, SG, USA). Format: --ipv6-policy "region_name:<addr1:str>,<addr2:str>"
    + tcp-healthcheck: Specify TCP configuration for healthcheck. Format: --tcp-healthcheck "tcp_port:<port:int>"
    + http-healthcheck: Specify HTTP configuration for healthcheck. Format: --http-healthcheck "http_port:<port:int>;url_path:<path:str>;ok_codes:<code1:int>,<code2:int>;vhost:<domain:str>;interval:<interval:int>"
    Example: ./bizfly dns create-record --zone-id zone_123 --name test_policy3 --ttl 600 --type GEOIP --ipv4-policy "HN:1.1.1.1,2.3.4.5" --ipv4-policy "HCM:68.68.79.23" --tcp-healthcheck "tcp_port:100" --http-healthcheck "http_port:1000;url_path:/health;ok_codes:400,404;vhost:test.com;interval:5"
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if checkValidType(recordType, NormalTypes) { // Normal type case
			stringRecordData := fmt.Sprintf("%v", recordData)
			data := parseNormalRecord(stringRecordData)
			if len(data) == 0 {
				log.Fatal("Invalid argument")
			}
			payloadData := gobizfly.CreateNormalRecordPayload{
				BaseCreateRecordPayload: gobizfly.BaseCreateRecordPayload{
					Name: recordName,
					Type: recordType,
					TTL:  TTL,
				},
				Data: data,
			}
			payload := recordPayload{
				Record: payloadData,
			}
			fmt.Printf("%+v\n", payload)
			json_data, _ := json.Marshal(payload)
			fmt.Println(string(json_data))
			recordSet, err := client.DNS.CreateRecord(ctx, zoneID, payload)
			if err != nil {
				log.Fatal(err)
			}
			var recordSetData [][]string
			stringData := ""
			for _, ip := range recordSet.Data {
				stringData += fmt.Sprintf("%v", ip) + "\n"
			}
			recordSetData = append(recordSetData, []string{recordSet.ID, recordSet.Name,
				recordSet.Type, strconv.Itoa(recordSet.TTL), stringData})
			formatter.Output(recordSetHeader, recordSetData)
			outputRecordData(recordSet)
		} else if recordType == "GEOIP" {
			if len(ipv4RoutingPolicy) == 0 && len(ipv6RoutingPolicy) == 0 {
				log.Fatal("Invalid argument")
			}
			v4Addrs := parsePolicyRecord(ipv4RoutingPolicy)
			v6Addrs := parsePolicyRecord(ipv6RoutingPolicy)
			HealthCheckPayload := gobizfly.HealthCheck{}
			if httpHealthCheck != "" {
				httpHealthCheckData := parseHTTPHealthCheck(httpHealthCheck)
				HealthCheckPayload.HTTPStatus = httpHealthCheckData
			}
			if tcpHealthCheck != "" {
				tcpHealthCheckData := parseTCPHealthCheck(tcpHealthCheck)
				HealthCheckPayload.TCPConnect = tcpHealthCheckData
			}
			payloadData := gobizfly.CreatePolicyRecordPayload{
				BaseCreateRecordPayload: gobizfly.BaseCreateRecordPayload{
					Name: recordName,
					Type: recordType,
					TTL:  TTL,
				},
				RoutingPolicyData: gobizfly.RoutingPolicyData{
					RoutingData: gobizfly.RoutingData{
						AddrsV4: v4Addrs,
						AddrsV6: v6Addrs,
					},
					HealthCheck: HealthCheckPayload,
				},
			}
			payload := recordPayload{
				Record: payloadData,
			}
			json_data, _ := json.Marshal(payload)
			fmt.Println(string(json_data))

			recordSet, err := client.DNS.CreateRecord(ctx, zoneID, payload)
			if err != nil {
				log.Fatal(err)
			}
			var recordSetData [][]string
			stringData := ""
			for _, ip := range recordSet.Data {
				stringData += fmt.Sprintf("%v", ip) + "\n"
			}
			recordSetData = append(recordSetData, []string{recordSet.ID, recordSet.Name,
				recordSet.Type, strconv.Itoa(recordSet.TTL), stringData})
			formatter.Output(recordSetHeader, recordSetData)
			outputRecordData(recordSet)
		} else if recordType == "MX" {
			mxData := parseMXRecord(domainData)
			payloadData := gobizfly.CreateMXRecordPayload{
				BaseCreateRecordPayload: gobizfly.BaseCreateRecordPayload{
					Name: recordName,
					TTL:  TTL,
					Type: recordType,
				},
				Data: mxData,
			}
			payload := recordPayload{
				Record: payloadData,
			}
			recordSet, err := client.DNS.CreateRecord(ctx, zoneID, payload)
			if err != nil {
				log.Fatal(err)
			}
			var recordSetData [][]string
			recordSetData = append(recordSetData, []string{recordSet.ID, recordSet.Name,
				recordSet.Type, strconv.Itoa(recordSet.TTL)})
			formatter.Output(recordSetHeader, recordSetData)
			outputRecordData(recordSet)
		}
	},
}

var deleteRecordCommand = &cobra.Command{
	Use:   "delete-record",
	Short: "Delete DNS record",
	Long: `Delete DNS record
Usage: ./bizfly dns delete-record <record-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		err := client.DNS.DeleteRecord(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Deleted record successfully")
	},
}

func outputRecordData(record *gobizfly.Record) {
	if checkValidType(record.Type, NormalTypes) {
		var IPs [][]string
		for _, ip := range record.Data {
			IPs = append(IPs, []string{ip.(string)})
		}
		formatter.Output(NormalDataHeader, IPs)
	} else if record.Type == "MX" {
		var mxDatas [][]string
		for _, domainData := range record.Data {
			domainMap := domainData.(map[string]interface{})
			priority := strconv.Itoa(int(domainMap["priority"].(float64)))
			mxDatas = append(mxDatas, []string{domainMap["value"].(string), priority})
		}
		formatter.Output(MXDataHeader, mxDatas)
	} else if record.Type == "GEOIP" {
		var RoutingData [][]string
		routingData := record.RoutingPolicyData.RoutingData
		RoutingData = append(RoutingData,
			[]string{"HN", joinIps(routingData.AddrsV4.HN), joinIps(routingData.AddrsV6.HN)})
		RoutingData = append(RoutingData,
			[]string{"HCM", joinIps(routingData.AddrsV4.HCM), joinIps(routingData.AddrsV6.HCM)})
		RoutingData = append(RoutingData,
			[]string{"SG", joinIps(routingData.AddrsV4.SG), joinIps(routingData.AddrsV6.SG)})
		RoutingData = append(RoutingData,
			[]string{"USA", joinIps(routingData.AddrsV4.USA), joinIps(routingData.AddrsV6.USA)})
		formatter.Output(PolicyRoutingHeader, RoutingData)
	}
}

func joinIps(data []string) string {
	return strings.Join(data, "\n")
}

func checkValidType(recordType string, validTypes []string) bool {
	var result bool
	for _, validType := range validTypes {
		if recordType == validType {
			result = true
		}
	}
	return result
}

func parseNormalRecord(data string) []string {
	return strings.Split(data, ";")
}

func parseMXRecord(data []string) []gobizfly.MXData {
	var mxData []gobizfly.MXData
	for _, recordString := range data {
		fragments := strings.Split(recordString, ":")
		domain := fragments[0]
		priority, err := strconv.Atoi(fragments[1])
		if err != nil {
			log.Fatal(err)
		}
		mxData = append(mxData, gobizfly.MXData{Value: domain, Priority: priority})
	}
	return mxData
}

func parsePolicyRecord(data []string) gobizfly.Addrs {
	addrs := gobizfly.Addrs{
		HN:  []string{},
		HCM: []string{},
		SG:  []string{},
		USA: []string{},
	}
	for _, regionData := range data {
		fragments := strings.Split(regionData, ":")
		region := fragments[0]
		ips := strings.Split(fragments[1], ",")
		switch region {
		case "HN":
			addrs.HN = ips
		case "HCM":
			addrs.HCM = ips
		case "SG":
			addrs.SG = ips
		case "USA":
			addrs.USA = ips
		}
	}
	return addrs
}

func parseHTTPHealthCheck(healthCheckData string) gobizfly.HTTPHealthCheck {
	var httpHealthCheck gobizfly.HTTPHealthCheck
	pairFields := strings.Split(healthCheckData, ";")
	for _, pairField := range pairFields {
		keyAndValue := strings.Split(pairField, ":")
		key := keyAndValue[0]
		value := keyAndValue[1]
		switch key {
		case "http_port":
			port, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
			httpHealthCheck.HTTPPort = port
		case "url_path":
			httpHealthCheck.URLPath = value
		case "vhost":
			httpHealthCheck.VHost = value
		case "ok_codes":
			var codes []int
			for _, code := range strings.Split(value, ",") {
				intCode, err := strconv.Atoi(code)
				if err != nil {
					log.Fatal(err)
				}
				codes = append(codes, intCode)
			}
			httpHealthCheck.OkCodes = codes
		case "interval":
			intInterval, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
			httpHealthCheck.Interval = intInterval
		}
	}
	return httpHealthCheck
}

func parseTCPHealthCheck(data string) gobizfly.TCPHealthCheck {
	var tcpHealthCheck gobizfly.TCPHealthCheck
	keyAndValue := strings.Split(data, ":")
	key := keyAndValue[0]
	value := keyAndValue[1]
	if key == "tcp_port" {
		intPort, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
		tcpHealthCheck.TCPPort = intPort
	}
	return tcpHealthCheck
}

func init() {
	rootCmd.AddCommand(dnsComnmand)
	dnsComnmand.AddCommand(listZonesCommand)

	dnsComnmand.AddCommand(getZoneCommand)

	czpf := createZoneCommand.PersistentFlags()
	czpf.BoolVar(&required, "required", false, "Is required")
	czpf.StringVar(&zoneDescription, "description", "", "Zone description")
	dnsComnmand.AddCommand(createZoneCommand)

	dnsComnmand.AddCommand(deleteZoneCommand)

	crpf := createRecordCommand.PersistentFlags()
	crpf.StringVar(&zoneID, "zone-id", "", "Zone ID")
	crpf.StringVar(&recordName, "name", "", "Name of record")
	crpf.StringVar(&recordType, "type", "", "Record type")
	crpf.IntVar(&TTL, "ttl", 0, "TTL of record")
	crpf.StringVar(&recordData, "data", "", "Data of record")
	crpf.StringArrayVar(&ipv4RoutingPolicy, "ipv4-policy", []string{}, "IPv4 policies")
	crpf.StringArrayVar(&ipv6RoutingPolicy, "ipv6-policy", []string{}, "IPv6 policies")
	crpf.StringVar(&tcpHealthCheck, "tcp-healthcheck", "", "TCP Health Check Configuration")
	crpf.StringVar(&httpHealthCheck, "http-healthcheck", "", "HTTP Health Check Configuration")
	crpf.StringArrayVar(&domainData, "domain-data", []string{}, "Domain with its priority")
	_ = cobra.MarkFlagRequired(crpf, "zone-id")
	_ = cobra.MarkFlagRequired(crpf, "name")
	_ = cobra.MarkFlagRequired(crpf, "type")
	_ = cobra.MarkFlagRequired(crpf, "ttl")
	dnsComnmand.AddCommand(createRecordCommand)

	dnsComnmand.AddCommand(deleteRecordCommand)

	dnsComnmand.AddCommand(getRecordCommand)
}
