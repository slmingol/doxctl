/*
Package cmd - ...

Copyright © 2021 Sam Mingolelli <github@lamolabs.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type svc struct {
	Svc      string   `mapstructure:"svc"`
	Svrs     []string `mapstructure:"svrs"`
	Port     int      `mapstructure:"port"`     // Optional: defaults to 6443
	Path     string   `mapstructure:"path"`     // Optional: defaults to /healthz
	Insecure bool     `mapstructure:"insecure"` // Optional: skip TLS verification for this service
}

type config struct {
	MinVpnRoutes     int           `mapstructure:"minVpnRoutes"`
	DomNameChk       string        `mapstructure:"domNameChk"`
	DomSearchChk     string        `mapstructure:"domSearchChk"`
	DomAddrChk       string        `mapstructure:"domAddrChk"`
	DomainName       string        `mapstructure:"domainName"`
	ServerA          string        `mapstructure:"digProbeServerA"`
	ServerB          string        `mapstructure:"digProbeServerB"`
	Sites            []string      `mapstructure:"sites"`
	Openshift        []string      `mapstructure:"openshift"`
	Svcs             []svc         `mapstructure:"wellKnownSvcs"`
	PingTimeout      time.Duration `mapstructure:"pingTimeout"`
	DNSLookupTimeout time.Duration `mapstructure:"dnsLookupTimeout"`
	FailThreshold    int           `mapstructure:"failThreshold"`
}

// Validate checks if the configuration is valid
func (c *config) Validate() error {
	if c.DomainName == "" {
		return fmt.Errorf("domainName is required in configuration file")
	}

	if len(c.Sites) == 0 {
		return fmt.Errorf("at least one site must be defined in the 'sites' configuration")
	}

	if len(c.Svcs) == 0 {
		return fmt.Errorf("at least one service must be defined in 'wellKnownSvcs' configuration")
	}

	// Validate each service has required fields
	for i, svc := range c.Svcs {
		if svc.Svc == "" {
			return fmt.Errorf("wellKnownSvcs[%d]: 'svc' field is required", i)
		}
		if len(svc.Svrs) == 0 {
			return fmt.Errorf("wellKnownSvcs[%d] (%s): at least one server must be defined in 'svrs'", i, svc.Svc)
		}
	}

	return nil
}

// setDefaults sets default values for optional configuration fields
func (c *config) setDefaults() {
	if c.PingTimeout == 0 {
		c.PingTimeout = 250 * time.Millisecond
	}
	if c.DNSLookupTimeout == 0 {
		c.DNSLookupTimeout = 100 * time.Millisecond
	}
	if c.FailThreshold == 0 {
		c.FailThreshold = 5
	}
	if c.MinVpnRoutes == 0 {
		c.MinVpnRoutes = 5
	}
}

var (
	cfgFile            string
	verboseChk, allChk bool
	conf               *config
	outputFormat       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "doxctl",
	Version: version,
	Short:   "A CLI to help triage network/DNS/VPN connectivity issues",
	Long: `
'doxctl' is a collection of tools which can be used to diagnose & triage problems 
stemming from the following areas with a laptop or desktop system:

  - DNS, specifically with the configuration of resolvers 
  - VPN configuration and network connectivity over it
  - General access to well-known servers
  - ... or general network connectivity issues 

	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())

	fmt.Printf("\n")
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.doxctl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseChk, "verbose", "v", false, "Enable verbose output of commands")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json, yaml")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Printf("error reading flags: %s\n", err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".doxctl" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".doxctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Build the message with proper padding to 80 chars
		configFile := viper.ConfigFileUsed()
		textContent := fmt.Sprintf(" ℹ NOTE  Using config file: %s", configFile)
		// Calculate padding needed (accounting for ANSI codes that won't be visible)
		padding := 80 - len(textContent)
		if padding < 0 {
			padding = 0
		}

		fmt.Println("")
		fmt.Printf("\033[38;2;0;128;128m%s\033[0m\n", strings.Repeat("─", 80))
		fmt.Printf("\033[48;2;0;64;64m\033[38;2;0;255;255;1m ℹ NOTE \033[0m\033[48;2;0;64;64m\033[1;97m Using config file: \033[38;2;135;206;250;1m%s\033[1;97m%s\033[0m\n",
			configFile, strings.Repeat(" ", padding))
		fmt.Printf("\033[38;2;0;128;128m%s\033[0m\n", strings.Repeat("─", 80))
		//fmt.Fprintln(os.Stderr, "\n**NOTE:** Using config file:", viper.ConfigFileUsed(), "\n")
	} else {
		// Config file is optional for some commands (like 'config init')
		// Commands that require config will validate in their PreRun
		return
	}

	conf = &config{}
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: Failed to parse configuration file '%s'\n", viper.ConfigFileUsed())
		fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Please check your configuration file for syntax errors.\n")
		fmt.Fprintf(os.Stderr, "See .doxctl.yaml.example for a sample configuration.\n\n")
		os.Exit(1)
	}

	// Set default values
	conf.setDefaults()

	// Validate configuration
	if err := conf.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: Invalid configuration in '%s'\n", viper.ConfigFileUsed())
		fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Please fix the configuration errors above.\n")
		fmt.Fprintf(os.Stderr, "See .doxctl.yaml.example for a sample configuration.\n\n")
		os.Exit(1)
	}
}
