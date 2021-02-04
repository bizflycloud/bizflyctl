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
	"bufio"
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	sshListHeader = []string{"Name", "Fingerprint"}
	sshKeyName    string
	publicKey     string
)

var sshKeyCmd = &cobra.Command{
	Use:   "ssh-key",
	Short: "BizFly Cloud SSH Key Interaction",
	Long:  `BizFly Cloud SSH Key Action: Create, List, Delete`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SSH Key called")
	},
}

var sshkeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your SSH keys",
	Long:  "List your SSH keys",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		keys, err := client.SSHKey.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, key := range keys {
			s := []string{key.SSHKeyPair.Name, key.SSHKeyPair.FingerPrint}
			data = append(data, s)
		}
		formatter.Output(sshListHeader, data)
	},
}

var sshKeyDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete your SSH key",
	Long:  "Delete a SSH Key using its name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Invalid arguments")
			os.Exit(1)
		}
		client, ctx := getApiClient(cmd)
		_, err := client.SSHKey.Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Deleted the SSH key")
	},
}

var sshKeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a SSH Key",
	Long: `Create a SSH Key using name and public key
Example 1: bizfly ssh-key create --name abcxyz --public-key your-pub-key
Example 2: bizfly ssh-key create --name abcxyz --public-key prompt-url => Then type your URL which contains your public key
Example 3: bizfly ssh-key create --name abcxyz --public-key prompt => Paste your public key, and then send EOF (Ctrl + D in *nix; Ctrl + Z in Windows)`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if publicKey == "prompt-url" {
			var key string
			fmt.Scanln(&key)
			fmt.Println(key)
			resp, err := http.Get(key)
			if err != nil {
				log.Fatal(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			publicKey = string(body)
		} else if publicKey == "prompt" {
			scanner := bufio.NewScanner(os.Stdin)
			line := scanner.Text()
			publicKey = line
		}
		key, err := client.SSHKey.Create(ctx, &gobizfly.SSHKeyCreateRequest{
			Name:      sshKeyName,
			PublicKey: publicKey,
		})
		if err != nil {
			log.Fatal(err)
		}
		data := [][]string{{key.Name, key.FingerPrint}}
		formatter.Output(sshListHeader, data)
	},
}

func init() {
	rootCmd.AddCommand(sshKeyCmd)
	sshKeyCmd.AddCommand(sshkeyListCmd)
	sshKeyCmd.AddCommand(sshKeyDeleteCmd)
	scpf := sshKeyCreateCmd.PersistentFlags()
	scpf.StringVar(&publicKey, "public-key", "", "The Public Key")
	scpf.StringVar(&sshKeyName, "name", "", "The SSH Key name")
	_ = cobra.MarkFlagRequired(scpf, "name")
	sshKeyCmd.AddCommand(sshKeyCreateCmd)
}
