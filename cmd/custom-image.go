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
	"fmt"
	"github.com/bizflycloud/bizflyctl/formatter"
	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	customImageHeader = []string{"ID", "Name", "Description", "Container Format", "Size", "Status", "Visibility"}
	imageURL          string
	diskFormat        string
	customImageName   string
	filePath          string
	downloadPath      string
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
		if imageURL == "" && filePath == "" {
			log.Fatal("Invalid arguments. You need to specify image-url or file-path")
		} else if imageURL != "" && filePath != "" {
			log.Fatal("Invalid arguments. You need to specify image-url or file-path")
		}
		client, ctx := getApiClient(cmd)
		if imageURL != "" {
			resp, err := client.Server.CreateCustomImage(ctx, &gobizfly.CreateCustomImagePayload{
				Name:        customImageName,
				DiskFormat:  diskFormat,
				Description: description,
				ImageURL:    imageURL,
			})
			if err != nil {
				log.Fatal(err)
			}
			image := resp.Image
			var data [][]string
			data = append(data, []string{image.ID, image.Name, image.Description,
				image.ContainerFormat, strconv.Itoa(image.Size), image.Status, image.Visibility})
			formatter.Output(customImageHeader, data)
		} else {
			resp, err := client.Server.CreateCustomImage(ctx, &gobizfly.CreateCustomImagePayload{
				Name:        customImageName,
				DiskFormat:  diskFormat,
				Description: description,
			})
			if err != nil {
				log.Fatal(err)
			}
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp.UploadURI)
			r, err := http.NewRequest("PUT", resp.UploadURI, file)

			if err != nil {
				log.Fatal(err)
			}
			r.Header.Set("X-Auth-Token", resp.Token)
			r.Header.Set("Content-Type", "application/octet-stream")
			client := &http.Client{}
			response, err := client.Do(r)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			image := resp.Image
			var data [][]string
			data = append(data, []string{image.ID, image.Name, image.Description,
				image.ContainerFormat, strconv.Itoa(image.Size), image.Status, image.Visibility})
			formatter.Output(customImageHeader, data)
		}
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
			fmt.Println("Delete the custom image successfully")
		}
	},
}

var customImageDownload = &cobra.Command{
	Use:   "download",
	Short: "Download a custom image",
	Long:  "Download a custom image using its ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := getApiClient(cmd)
		resp, err := client.Server.GetCustomImage(ctx, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		image := resp.Image
		token := resp.Token

		if image.ID == args[0] {
			downloadURL := image.File
			fileName := fmt.Sprintf("%s.%s", image.Name, image.ContainerFormat)
			file, err := os.Create(filepath.Join(downloadPath, fileName))
			if err != nil {
				log.Fatal(err)
			}
			client := http.Client{}
			req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("X-Auth-Token", token)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			size, err := io.Copy(file, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			fmt.Printf("Downloaded a file %s with size %d Bytes\n", fileName, size)

			data = append(data, []string{image.ID, image.Name, image.Description,
				image.ContainerFormat, strconv.Itoa(image.Size), image.Status, image.Visibility})
		}
		formatter.Output(customImageHeader, data)
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
	ccpf.StringVar(&filePath, "file-path", "", "Upload file path")
	_ = cobra.MarkFlagRequired(ccpf, "name")
	_ = cobra.MarkFlagRequired(ccpf, "disk-format")

	customImageCmd.AddCommand(customImageCreate)
	customImageCmd.AddCommand(customImageDownload)
	cdpf := customImageDownload.PersistentFlags()
	cdpf.StringVar(&downloadPath, "output-path", ".", "Output file path")
}
