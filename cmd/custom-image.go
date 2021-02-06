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
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var (
	customImageHeader = []string{"ID", "Name", "Description", "ContainerFormat", "Size", "Status", "Visibility"}
	customImageID     string
	imageURL          string
	diskFormat        string
	customImageName   string
)

var customImageCmd = &cobra.Command{
	Use:   "custom-image",
	Short: "BizFly Custom Image Interaction",
	Long:  "BizFly Custom Image Action: List, Create, Delete",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("custom image called")
	},
}

var customImageList = &cobra.Command{
	Use:   "list",
	Short: "List custom images",
	Long:  "List your custom images",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		images, err := client.Server.ListCustomImages(ctx)
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, image := range images {
			data = append(data, []string{image.ID, image.Name, image.Description,
				image.ContainerFormat, strconv.Itoa(image.Size), image.Status, image.Visibility})
		}
		formatter.Output(customImageHeader, data)
	},
}

var customImageCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new custom image",
	Long: `Create a new custom image with name, image URL
Example: bizfly custom-image create --name xyz --disk-format raw --description abcxyz --image-url http://xyz.abc`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		image, err := client.Server.CreateCustomImage(ctx, &gobizfly.CreateCustomImagePayload{
			Name:        customImageName,
			DiskFormat:  diskFormat,
			Description: description,
			ImageURL:    imageURL,
		})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{image.ID, image.Name, image.Description,
			image.ContainerFormat, strconv.Itoa(image.Size), image.Status, image.Visibility})
		formatter.Output(customImageHeader, data)
	},
}

var customImageDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a custom image",
	Long:  "Delete a custom image using its ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Invalid argument")
		}
		client, ctx := getApiClient(cmd)
		err := client.Server.DeleteCustomImage(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Delete the custom image sucessfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(customImageCmd)
	customImageCmd.AddCommand(customImageList)
	customImageCmd.AddCommand(customImageDelete)

	ccpf := customImageCreate.PersistentFlags()
	ccpf.StringVar(&customImageName, "name", "", "Name of custom image")
	ccpf.StringVar(&diskFormat, "disk-format", "", "Disk format of image")
	ccpf.StringVar(&description, "description", "", "Description")
	ccpf.StringVar(&imageURL, "image-url", "", "Image URL")
	_ = cobra.MarkFlagRequired(ccpf, "name")
	_ = cobra.MarkFlagRequired(ccpf, "disk-format")
	customImageCmd.AddCommand(customImageCreate)
}
