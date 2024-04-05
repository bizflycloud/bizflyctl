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
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bizflycloud/bizflyctl/constants"
	"github.com/bizflycloud/gobizfly"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	email         string
	password      string
	region        string
	project_id    string
	appCredSecret string
	appCredID     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bizfly",
	Short: "Bizfly Cloud Command Line",
	Long:  `Bizfly Cloud Command Line`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pre run")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bizfly.yaml)")

	rootCmd.PersistentFlags().StringVar(&email, "email", "", "Your Bizfly Cloud Email. Read environment variable BIZFLY_CLOUD_EMAIL")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Your Bizfly Cloud Password. Read environment variable BIZFLY_CLOUD_PASSWORD")
	rootCmd.PersistentFlags().StringVar(&appCredID, "app-credential-id", "", "Your Bizfly Cloud Application Credential Id. Read environment variable BIZFLY_CLOUD_APP_CREDENTIAL_ID")
	rootCmd.PersistentFlags().StringVar(&appCredSecret, "app-credential-secret", "", "Your Bizfly Cloud Application Credential Secret. Read environment variable BIZFLY_CLOUD_APP_CREDENTIAL_SECRET")

	rootCmd.PersistentFlags().StringVar(&region, "region", "HaNoi", "Region you want to access the resource. Read environment variable BIZFLY_CLOUD_REGION")
	rootCmd.PersistentFlags().StringVar(&project_id, "project-id", "", "Your Bizfly Cloud Project ID. Read environment variable BIZFLY_CLOUD_PROJECT_ID")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bizfly" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bizfly")
	}

	viper.SetEnvPrefix("BIZFLY_CLOUD")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getRegionName(regionName string) string {
	lowerRegion := strings.ToLower(regionName)
	result := constants.RegionMapping[lowerRegion]
	return result
}

func getApiClient(cmd *cobra.Command) (*gobizfly.Client, context.Context) {
	// use application credential auth
	if appCredID == "" {
		appCredID = viper.GetString("app_credential_id")
	}
	if appCredSecret == "" {
		appCredSecret = viper.GetString("app_credential_secret")
	}
	useAppCredential := true
	if appCredID == "" && appCredSecret == "" {
		// use username/password auth
		if email == "" {
			email = viper.GetString("email")
		}
		if password == "" {
			password = viper.GetString("password")
		}
		useAppCredential = false
	}

	if viper.GetString("region") != "" {
		region = viper.GetString("region")
	}

	regionName := getRegionName(region)
	if regionName == "" {
		log.Fatalf("Invalid region %s", region)
	}

	if viper.GetString("project_id") != "" {
		project_id = viper.GetString("project_id")
	}
	// nolint:staticcheck
	client, err := gobizfly.NewClient(gobizfly.WithProjectId(project_id), gobizfly.WithRegionName(regionName))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	//TODO Get token from cache
	request := &gobizfly.TokenCreateRequest{
		ProjectID: project_id,
	}
	if useAppCredential {
		request.AuthType = gobizfly.AppCredentialAuthType
		request.AppCredID = appCredID
		request.AppCredSecret = appCredSecret
	} else {
		request.AuthMethod = "password"
		request.Username = email
		request.Password = password
	}
	tok, err := client.Token.Create(ctx, request)
	if err != nil {
		log.Fatal(err)
	}
	client.SetKeystoneToken(tok)
	ctx = context.WithValue(ctx, "token", tok.KeystoneToken)
	return client, ctx
}
