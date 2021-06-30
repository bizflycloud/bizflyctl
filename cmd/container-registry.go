/*
Copyright © 2021 izFly Cloud

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
)

var (
	isPublic bool
	isPrivate bool
	repoName string
	tagName string
	vulnerabilities string
	repositoryHeader = []string{"Name", "LastPush", "Pulls", "Public", "CreatedAt"}
	vulnerabilityHeader = []string{"Package", "Name", "Namespace", "Link", "Severity", "FixedBy"}
	tagHeader = []string{"Name", "Author", "LastUpdated", "CreatedAt", "LastScan", "ScanStatus", "Vulnerabilities", "Fixes"}
)

var containerRegistryCmd = &cobra.Command{
	Use: "container-registry",
	Short: "BizFly Cloud Container Registry Interaction",
	Long: "BizFly CLoud Container Registry Action: List, Create, Delete, Get Tags, Update, Delete Image Tag, Get Image Info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("container registry called")
	},
}

var repositoryListCmd = &cobra.Command{
	Use: "list",
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
	Use: "create",
	Short: "Create Container Registry repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		payload := &gobizfly.CreateRepositoryPayload{
			Name: repoName,
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
	Use: "delete",
	Short: "Delete Container Registry repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		err := client.ContainerRegistry.Delete(ctx, repoName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Deleting repository")
	},
}

var getTagCmd = &cobra.Command{
	Use: "get-tags",
	Short: "Get repository Tags",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		repoTags, err := client.ContainerRegistry.GetTags(ctx, repoName)
		if err != nil {
			log.Fatal(err)
		}
		repo := repoTags.Repository
		var repoData [][]string
		repoData = append(repoData, []string{repo.Name, repo.LastPush, strconv.Itoa(repo.Pulls), strconv.FormatBool(repo.Public), repo.CreatedAt})
		formatter.Output(repositoryHeader, repoData)

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
	Use: "edit-repo",
	Short: "Edit Container Registry repository",
	Run: func(cmd *cobra.Command, args []string) {
		if !isPublic && !isPrivate {  // Both two vars don't set
			log.Fatal("You need to choose is-public or is-private")
		}
		isPublic = isPublic || !isPrivate
		client, ctx := getApiClient(cmd)
		payload := &gobizfly.EditRepositoryPayload{
			Public: isPublic,
		}
		err := client.ContainerRegistry.EditRepo(ctx, repoName, payload)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Edit repository successfully")
	},
}

var deleteTagCmd = &cobra.Command{
	Use: "delete-tag",
	Short: "Delete Repository Tag",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		err := client.ContainerRegistry.DeleteTag(ctx, tagName, repoName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Delete tag of repository successfully")
	},
}

var getImageCmd = &cobra.Command{
	Use: "get-image",
	Short: "Get repository tag",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		image, err := client.ContainerRegistry.GetTag(ctx, repoName, tagName, vulnerabilities)
		if err != nil {
			log.Fatal(err)
		}
		var repoData [][]string
		repo := image.Repository
		repoData = append(repoData, []string{repo.Name, repo.LastPush, strconv.Itoa(repo.Pulls), strconv.FormatBool(repo.Public), repo.CreatedAt})
		formatter.Output(repositoryHeader, repoData)
		var tagsData [][]string

		tag := image.Tag
		tagsData = append(tagsData, []string{tag.Name, tag.Author, tag.LastUpdated, tag.CreatedAt, tag.LastScan,
			tag.ScanStatus, strconv.Itoa(tag.Vulnerabilities), strconv.Itoa(tag.Fixes)})
		formatter.Output(tagHeader, tagsData)

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

func init() {
	rootCmd.AddCommand(containerRegistryCmd)

	containerRegistryCmd.AddCommand(repositoryListCmd)

	rcpf := repositoryCreateCmd.PersistentFlags()
	rcpf.StringVar(&repoName, "repo-name", "", "Repository name")
	rcpf.BoolVar(&isPublic, "is-public", false, "Is public or not")
	_ = cobra.MarkFlagRequired(rcpf, "repo-name")
	containerRegistryCmd.AddCommand(repositoryCreateCmd)

	rdpf := repositoryDeleteCmd.PersistentFlags()
	rdpf.StringVar(&repoName, "repo-name", "", "Repository name")
	_ = cobra.MarkFlagRequired(rdpf, "repo-name")
	containerRegistryCmd.AddCommand(repositoryDeleteCmd)

	gtpf := getTagCmd.PersistentFlags()
	gtpf.StringVar(&repoName, "repo-name", "", "Repository name")
	_ = cobra.MarkFlagRequired(gtpf, "repo-name")
	containerRegistryCmd.AddCommand(getTagCmd)

	erpf := editRepoCmd.PersistentFlags()
	erpf.StringVar(&repoName, "repo-name", "", "Repository name")
	erpf.BoolVar(&isPublic, "is-public", false, "Is public or not")
	erpf.BoolVar(&isPrivate, "is-private", false, "Is private or not")
	_ = cobra.MarkFlagRequired(erpf, "repo-name")
	containerRegistryCmd.AddCommand(editRepoCmd)

	dtpf := deleteTagCmd.PersistentFlags()
	dtpf.StringVar(&repoName, "repo-name", "", "Repository name")
	dtpf.StringVar(&tagName, "tag-name", "", "Tag name")
	_ = cobra.MarkFlagRequired(dtpf, "repo-name")
	_ = cobra.MarkFlagRequired(dtpf, "tag-name")
	containerRegistryCmd.AddCommand(deleteTagCmd)

	gipf := getImageCmd.PersistentFlags()
	gipf.StringVar(&repoName, "repo-name", "", "Repository name")
	gipf.StringVar(&tagName, "tag-name", "", "Tag name")
	gipf.StringVar(&vulnerabilities, "vulnerabilities", "", "Vulnerabilities")
	_ = cobra.MarkFlagRequired(gipf, "repo-name")
	_ = cobra.MarkFlagRequired(gipf, "tag-name")
	containerRegistryCmd.AddCommand(getImageCmd)


}