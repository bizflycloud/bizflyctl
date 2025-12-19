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
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Bizfly Cloud via browser",
	Long:  `Login to Bizfly Cloud via browser to obtain an authentication token.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLogin(); err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin() error {
	// 1. Start a local HTTP server on port 15995
	listener, err := net.Listen("tcp", "localhost:15995")
	if err != nil {
		return fmt.Errorf("failed to start local server: %w", err)
	}
	defer listener.Close()

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
		// Get ticket from query params
		ticket := r.URL.Query().Get("ticket")
		if ticket == "" {
			fmt.Fprintf(w, "Login failed: No ticket found in callback request. Query params: %v", r.URL.Query())
			errChan <- fmt.Errorf("no ticket found in callback")
			return
		}

		// Call serviceValidate endpoint to get real token
		validateURL := fmt.Sprintf("https://manage.bizflycloud.vn/cas/serviceValidate?ticket=%s&service=%s", ticket, callbackURL)
		resp, err := http.Get(validateURL)
		if err != nil {
			fmt.Fprintf(w, "Login failed: Failed to validate ticket: %v", err)
			errChan <- fmt.Errorf("failed to validate ticket: %w", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "Login failed: Failed to read validation response: %v", err)
			errChan <- fmt.Errorf("failed to read validation response: %w", err)
			return
		}

		// Parse XML response
		var casResponse CASServiceResponse
		if err := xml.Unmarshal(body, &casResponse); err != nil {
			fmt.Fprintf(w, "Login failed: Failed to parse validation response: %v", err)
			errChan <- fmt.Errorf("failed to parse validation response: %w", err)
			return
		}

		// Extract token from response
		if casResponse.AuthenticationSuccess == nil {
			fmt.Fprintf(w, "Login failed: Authentication failed. Response: %s", string(body))
			errChan <- fmt.Errorf("authentication failed")
			return
		}

		// Get token from attributes section
		if casResponse.AuthenticationSuccess.Attributes == nil {
			fmt.Fprintf(w, "Login failed: No attributes in validation response. Response: %s", string(body))
			errChan <- fmt.Errorf("no attributes in validation response")
			return
		}

		token := casResponse.AuthenticationSuccess.Attributes.Token
		if token == "" {
			fmt.Fprintf(w, "Login failed: No token in validation response. Response: %s", string(body))
			errChan <- fmt.Errorf("no token in validation response")
			return
		}

		fmt.Fprintf(w, "Login successful! You can close this window.")
		tokenChan <- token
	})

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case token := <-tokenChan:
		// 5. Save the token
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
		fmt.Println("Login successful! Token saved to config file.")
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
	XMLName              xml.Name              `xml:"http://www.yale.edu/tp/cas serviceResponse"`
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
