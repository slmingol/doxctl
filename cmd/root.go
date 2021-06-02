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
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type svc struct {
	Svc  string   `mapstructure:"svc"`
	Svrs []string `mapstructure:"svrs"`
}

type config struct {
	MinVpnRoutes  int           `mapstructure:"minVpnRoutes"`
	DomNameChk    string        `mapstructure:"domNameChk"`
	DomSearchChk  string        `mapstructure:"domSearchChk"`
	DomAddrChk    string        `mapstructure:"domAddrChk"`
	DomainName    string        `mapstructure:"domainName"`
	ServerA       string        `mapstructure:"digProbeServerA"`
	ServerB       string        `mapstructure:"digProbeServerB"`
	Sites         []string      `mapstructure:"sites"`
	Openshift     []string      `mapstructure:"openshift"`
	Svcs          []svc         `mapstructure:"wellKnownSvcs"`
	PingTimeout   time.Duration `mapstructure:"pingTimeout"`
	FailThreshold int           `mapstructure:"failThreshold"`
}

var (
	cfgFile            string
	verboseChk, allChk bool
	conf               *config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "doxctl",
	Short: "A CLI to help triage network/DNS/VPN connectivity issues",
	Long: `
'doxctl' is a collection of tools which can be used to diagnose & triage problems 
stemming from the following areas with a laptop or desktop system:

  - DNS, specifically with the configuration of resolvers 
  - VPN configuration and network connectivity over it
  - General access to well-known servers
  - General access to well-known services
  - ... or general network connectivity issues 

	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())

	fmt.Println("\n\n")
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.doxctl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseChk, "verbose", "v", false, "Enable verbose output of commands")

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
		viper.SetConfigName(".doxctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("")
		color.Note.Tips("Using config file: " + viper.ConfigFileUsed() + "\n")
		//fmt.Fprintln(os.Stderr, "\n**NOTE:** Using config file:", viper.ConfigFileUsed(), "\n")
	}

	conf := &config{}
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
}
