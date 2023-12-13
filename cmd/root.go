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
	"time"

	"github.com/bizflycloud/gobizfly"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	email      string
	password   string
	region     string
	project_id string
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
	_ = rootCmd.MarkFlagRequired("email")

	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Your Bizfly Cloud Password. Read environment variable BIZFLY_CLOUD_PASSWORD")
	_ = rootCmd.MarkFlagRequired("password")

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

func getApiClient(cmd *cobra.Command) (*gobizfly.Client, context.Context) {

	if email == "" {
		email = viper.GetString("email")
	}
	if email == "" {
		log.Fatal("Email is required")
	}

	if password == "" {
		password = viper.GetString("password")
	}
	if password == "" {
		log.Fatal("Password is required")
	}

	if viper.GetString("region") != "" {
		region = viper.GetString("region")
		if region == "HN" || region == "HCM" {
			log.Println("HN and HCM is deprecated. Using HaNoi and HoChiMinh instead")
		}
	}

	if viper.GetString("project_id") != "" {
		project_id = viper.GetString("project_id")
	}
	// nolint:staticcheck
	client, err := gobizfly.NewClient(gobizfly.WithProjectId(project_id), gobizfly.WithRegionName(region))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	// TODO Get token from cache
	tok, err := client.Token.Create(ctx,
		&gobizfly.TokenCreateRequest{
			AuthMethod: "password",
			Username:   email,
			Password:   password,
			ProjectID:  project_id,
		})
	if err != nil {
		log.Fatal(err)
	}

	client.SetKeystoneToken(tok)
	ctx = context.WithValue(ctx, "token", tok.KeystoneToken)

	return client, ctx
}
