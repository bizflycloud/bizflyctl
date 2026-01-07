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
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Bizfly Cloud via browser",
	Long:  `Login to Bizfly Cloud via browser to obtain an authentication token.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLogin(cmd); err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command) error {
	// 1. Start a local HTTP server on port 15995
	listener, err := net.Listen("tcp", "localhost:15995")
	if err != nil {
		return fmt.Errorf("failed to start local server: %w", err)
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Printf("failed to close login listener: %v", cerr)
		}
	}()

	port := 15995
	callbackURL := fmt.Sprintf("http://localhost:%d/callback", port)

	// 2. Construct the login URL
	loginURL := fmt.Sprintf("https://id.bizflycloud.vn/login?service=%s", callbackURL)

	fmt.Printf("Opening browser to login: %s\n", loginURL)

	// 3. Open the browser
	if err := openBrowser(loginURL); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
		fmt.Printf("Please open the URL manually.\n")
	}

	// 4. Wait for the callback
	tokenChan := make(chan string)
	errChan := make(chan error)

	server := &http.Server{}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		writeResponse := func(format string, args ...interface{}) {
			if _, err := fmt.Fprintf(w, format, args...); err != nil {
				log.Printf("failed to write login response: %v", err)
			}
		}

		// Get ticket from query params
		ticket := r.URL.Query().Get("ticket")
		if ticket == "" {
			writeResponse("Login failed: No ticket found in callback request. Query params: %v", r.URL.Query())
			errChan <- fmt.Errorf("no ticket found in callback")
			return
		}

		// Call serviceValidate endpoint to get real token
		validateURL := fmt.Sprintf("https://manage.bizflycloud.vn/cas/serviceValidate?ticket=%s&service=%s", ticket, callbackURL)
		resp, err := http.Get(validateURL)
		if err != nil {
			writeResponse("Login failed: Failed to validate ticket: %v", err)
			errChan <- fmt.Errorf("failed to validate ticket: %w", err)
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				log.Printf("failed to close validation response body: %v", cerr)
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			writeResponse("Login failed: Failed to read validation response: %v", err)
			errChan <- fmt.Errorf("failed to read validation response: %w", err)
			return
		}

		// Parse XML response
		var casResponse CASServiceResponse
		if err := xml.Unmarshal(body, &casResponse); err != nil {
			writeResponse("Login failed: Failed to parse validation response: %v", err)
			errChan <- fmt.Errorf("failed to parse validation response: %w", err)
			return
		}

		// Extract token from response
		if casResponse.AuthenticationSuccess == nil {
			writeResponse("Login failed: Authentication failed. Response: %s", string(body))
			errChan <- fmt.Errorf("authentication failed")
			return
		}

		// Get token from attributes section
		if casResponse.AuthenticationSuccess.Attributes == nil {
			writeResponse("Login failed: No attributes in validation response. Response: %s", string(body))
			errChan <- fmt.Errorf("no attributes in validation response")
			return
		}

		token := casResponse.AuthenticationSuccess.Attributes.Token
		if token == "" {
			writeResponse("Login failed: No token in validation response. Response: %s", string(body))
			errChan <- fmt.Errorf("no token in validation response")
			return
		}

		writeResponse("Login successful! You can close this window.")
		tokenChan <- token
	})

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case token := <-tokenChan:
		// 5. Check if project_id is provided
		// Get project_id from flag, config, or environment variable
		projID := ""
		// First try to get from persistent flag (check root command)
		if rootFlag := cmd.Root().PersistentFlags().Lookup("project-id"); rootFlag != nil {
			projID = rootFlag.Value.String()
		}
		// Fall back to package variable (set by persistent flag)
		if projID == "" {
			projID = project_id
		}
		// Finally check config/environment
		if projID == "" {
			projID = viper.GetString("project_id")
		}

		// If project_id is provided, exchange root token for project-scoped token
		if projID != "" {
			// Get region for client creation
			reg := region
			if reg == "" {
				reg = viper.GetString("region")
			}
			if reg == "" {
				reg = "HaNoi" // default region
			}

			regionName := getRegionName(reg)
			if regionName == "" {
				return fmt.Errorf("invalid region %s", reg)
			}

			// Exchange root token for project-scoped token
			// Make direct HTTP call since gobizfly library bypasses HTTP request when Token is set
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
			defer cancelFunc()

			// Construct the token endpoint URL based on region
			// The endpoint format is typically: https://manage.bizflycloud.vn/api/token
			tokenURL := "https://manage.bizflycloud.vn/api/token"

			// Prepare the request payload
			payload := map[string]string{
				"auth_method": "token",
				"token":       token,
				"project_id":  projID,
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to marshal request payload: %w", err)
			}

			// Create HTTP request
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBuffer(jsonPayload))
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			req.Header.Set("Content-Type", "application/json")

			// Make the HTTP request
			httpClient := &http.Client{
				Timeout: time.Second * 10,
			}
			resp, err := httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to exchange token: %w", err)
			}
			defer func() {
				if cerr := resp.Body.Close(); cerr != nil {
					log.Printf("failed to close token exchange response body: %v", cerr)
				}
			}()

			// Check response status
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
			}

			// Parse response
			var tokenResponse gobizfly.Token
			if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}

			// Use the project-scoped token
			if tokenResponse.KeystoneToken != "" {
				// Check if the token actually changed
				if tokenResponse.KeystoneToken == token {
					fmt.Printf("Warning: Token did not change after exchange. The API returned the same token.\n")
				}
				token = tokenResponse.KeystoneToken
			} else {
				return fmt.Errorf("received empty token from exchange")
			}
		}
		// Save the token
		viper.Set("auth_token", token)

		// Get the config file path
		configPath := viper.ConfigFileUsed()
		if configPath == "" {
			// If no config file is set, determine the default path
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home dir: %w", err)
			}
			configPath = fmt.Sprintf("%s/.bizfly.yaml", home)
			viper.SetConfigFile(configPath)
		}

		// Set config type to yaml to ensure viper knows the format
		viper.SetConfigType("yaml")

		// Check if config file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// File doesn't exist, use WriteConfigAs to create it with the full path
			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			// File exists, use WriteConfig to update it
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		if projID != "" {
			fmt.Println("Login successful! Project-scoped token saved to config file.")
		} else {
			fmt.Println("Login successful! Token saved to config file.")
		}
		return server.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// CASServiceResponse represents the CAS serviceValidate XML response
// XML namespace: http://www.yale.edu/tp/cas
type CASServiceResponse struct {
	XMLName               xml.Name               `xml:"http://www.yale.edu/tp/cas serviceResponse"`
	AuthenticationSuccess *AuthenticationSuccess `xml:"http://www.yale.edu/tp/cas authenticationSuccess"`
	AuthenticationFailure *AuthenticationFailure `xml:"http://www.yale.edu/tp/cas authenticationFailure"`
}

type AuthenticationSuccess struct {
	XMLName    xml.Name    `xml:"http://www.yale.edu/tp/cas authenticationSuccess"`
	User       string      `xml:"http://www.yale.edu/tp/cas user"`
	Attributes *Attributes `xml:"http://www.yale.edu/tp/cas attributes"`
}

type Attributes struct {
	Token string `xml:"http://www.yale.edu/tp/cas token"`
}

type AuthenticationFailure struct {
	XMLName xml.Name `xml:"http://www.yale.edu/tp/cas authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",chardata"`
}
