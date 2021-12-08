/*
Copyright © (2020-2021) izFly Cloud

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.rg/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

var (
	isPublic            bool
	isPrivate           bool
	expiresIn           int
	vulnerabilities     string
	scope               []string
	repositoryHeader    = []string{"Name", "Last Push", "Pulls", "Public", "Created At"}
	vulnerabilityHeader = []string{"Package", "Name", "Namespace", "Link", "Severity", "Fixed By"}
	tagHeader           = []string{"Name", "Author", "Last Updated", "Created At", "Last Scan", "Scan Status", "Vulnerabilities", "Fixes"}
)

var containerRegistryCmd = &cobra.Command{
	Use:   "container-registry",
	Short: "Bizfly Cloud Container Registry Interaction",
	Long:  "Bizfly Cloud Container Registry Action: List, Create, Delete, Get Tags, Update, Delete Image Tag, Get Image Info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("container registry called")
	},
}

var repositoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		repos, err := client.ContainerRegistry.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, repo := range repos {
			data = append(data, []string{repo.Name, repo.LastPush, strconv.Itoa(repo.Pulls), strconv.FormatBool(repo.Public), repo.CreatedAt})
		}
		formatter.Output(repositoryHeader, data)
	},
}

var repositoryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Container Registry repository",
	Long: `Create Container Registry repository
Usage: ./bizfly container-registry create <repo_name> (--public|--private)`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		if (!isPrivate && !isPublic) || (isPrivate && isPublic) {
			log.Fatal("You need to specify repository is public or not")
		}
		isPublic = isPublic || !isPrivate
		payload := &gobizfly.CreateRepositoryPayload{
			Name:   args[0],
			Public: isPublic,
		}
		err := client.ContainerRegistry.Create(ctx, payload)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Creating repository")
	},
}

var repositoryDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Container Registry repository",
	Long: `Delete Container Registry repository
Usage: ./bizfly container-registry delete <repo_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		err := client.ContainerRegistry.Delete(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Deleting repository")
	},
}

var getTagCmd = &cobra.Command{
	Use:   "get-tags",
	Short: "Get repository Tags",
	Long: `Get Repository Tags
Usage: ./bizfly container-registry get-tags <repo_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		repoTags, err := client.ContainerRegistry.GetTags(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var tagsData [][]string
		tags := repoTags.Tags
		for _, tag := range tags {
			tagsData = append(tagsData, []string{tag.Name, tag.Author, tag.LastUpdated, tag.CreatedAt, tag.LastScan,
				tag.ScanStatus, strconv.Itoa(tag.Vulnerabilities), strconv.Itoa(tag.Fixes)})
		}
		formatter.Output(tagHeader, tagsData)
	},
}

var editRepoCmd = &cobra.Command{
	Use:   "edit-repo",
	Short: "Edit Container Registry repository",
	Long: `Edit Container Registry repository
Usage: ./bizfly edit-repo <repo_name> (--public|--private)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		if (!isPrivate && !isPublic) || (isPrivate && isPublic) {
			log.Fatal("You need to specify repository is public or not")
		}
		isPublic = isPublic || !isPrivate
		client, ctx := getApiClient(cmd)
		payload := &gobizfly.EditRepositoryPayload{
			Public: isPublic,
		}
		err := client.ContainerRegistry.EditRepo(ctx, args[0], payload)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Edit repository successfully")
	},
}

var deleteTagCmd = &cobra.Command{
	Use:   "delete-tag",
	Short: "Delete Repository Tag",
	Long: `Delete Repository Tag
Usage: ./bizfly container-registry delete-tag <repo_name> <tag_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 2 {
			log.Fatal("Invalid argument")
		}
		err := client.ContainerRegistry.DeleteTag(ctx, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Delete tag of repository successfully")
	},
}

var getImageCmd = &cobra.Command{
	Use:   "get-image",
	Short: "Get repository tag",
	Long: `Get repository tag
Usage: ./bizfly container-registry get-image <repo_name> <tag_name> [flags]`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		if len(args) != 2 {
			log.Fatal("Invalid arguments")
		}
		image, err := client.ContainerRegistry.GetTag(ctx, args[0], args[1], vulnerabilities)
		if err != nil {
			log.Fatal(err)
		}
		vulnerabilities := image.Vulnerabilities
		var vulnerabilitiesData [][]string
		for _, vulnerability := range vulnerabilities {
			vulnerabilitiesData = append(vulnerabilitiesData,
				[]string{vulnerability.Package, vulnerability.Name, vulnerability.Namespace,
					vulnerability.Link, vulnerability.Severity, vulnerability.FixedBy})
		}
		formatter.Output(vulnerabilityHeader, vulnerabilitiesData)
	},
}

var genTokenCmd = &cobra.Command{
	Use:   "gen-token",
	Short: "Generate token for Container Registry",
	Long: `Generate token for Container Registry
Define:
- expires-in: Token expiration time in seconds, min: 1, max: 604800
- scope: Scopes which token grant to
   - actions: Action grant to the token [pull|push]
   - repository: Repository name or namespace (which ends with /). Leave blank in order to grant token to all repositories
Example: ./bizfly container-registry gen-token --expires-in 3404 --scope "actions:pull,push;repository:" --scope "actions:push;repository:test"
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		payload := &gobizfly.GenerateTokenPayload{
			ExpiresIn: expiresIn,
			Scopes:    parseScope(scope),
		}
		resp, err := client.ContainerRegistry.GenerateToken(ctx, payload)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Token:", resp.Token)
	},
}

func parseScope(scopes []string) []gobizfly.Scope {
	var scopeObjs []gobizfly.Scope
	for _, scope := range scopes {
		var scopeObj gobizfly.Scope
		fragments := strings.Split(scope, ";")
		if len(fragments) == 0 {
			log.Fatal("Invalid argument: scope")
		}
		for _, fragment := range fragments {
			keyValue := strings.Split(fragment, ":")
			if len(keyValue) != 2 {
				log.Fatal("Invalid argument: scope")
			}
			key := keyValue[0]
			value := keyValue[1]
			switch key {
			case "actions":
				actions := strings.Split(value, ",")
				scopeObj.Action = actions
			case "repository":
				scopeObj.Repository = value
			}
		}
		scopeObjs = append(scopeObjs, scopeObj)
	}
	return scopeObjs
}

func init() {
	rootCmd.AddCommand(containerRegistryCmd)

	containerRegistryCmd.AddCommand(repositoryListCmd)

	rcpf := repositoryCreateCmd.PersistentFlags()
	rcpf.BoolVar(&isPublic, "public", false, "Is public or not")
	rcpf.BoolVar(&isPrivate, "private", false, "Is private or not")
	containerRegistryCmd.AddCommand(repositoryCreateCmd)

	containerRegistryCmd.AddCommand(repositoryDeleteCmd)

	containerRegistryCmd.AddCommand(getTagCmd)

	erpf := editRepoCmd.PersistentFlags()
	erpf.BoolVar(&isPublic, "public", false, "Is public or not")
	erpf.BoolVar(&isPrivate, "private", false, "Is private or not")
	containerRegistryCmd.AddCommand(editRepoCmd)

	containerRegistryCmd.AddCommand(deleteTagCmd)

	getImageCmd.PersistentFlags().StringVar(&vulnerabilities, "vulnerabilities", "", "Image vulnerabilities")
	containerRegistryCmd.AddCommand(getImageCmd)

	gtopf := genTokenCmd.PersistentFlags()
	gtopf.IntVar(&expiresIn, "expires-in", 0, "Expires in (seconds)")
	gtopf.StringArrayVar(&scope, "scope", []string{}, "Token Scopes")
	_ = cobra.MarkFlagRequired(gtopf, "expires-in")
	_ = cobra.MarkFlagRequired(gtopf, "scope")
	containerRegistryCmd.AddCommand(genTokenCmd)
}
