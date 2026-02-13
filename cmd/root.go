/*
Copyright © 2026 Amanda Hager Lopes de Andrade Katz amandahla@gmail.com

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
	"os"

	"github.com/amandahla/syncli/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var timeout int
var config internal.Config
var debug bool

var logger = logrus.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syncli",
	Short: "CLI for interacting with Synapse Matrix homeserver",
	Long:  `SynCLI is a command-line interface for interacting with Synapse Matrix homeserver. It provides various commands to manage users, rooms, and other aspects of the Synapse server.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.syncli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "", 30, "Timeout for requests in seconds. (default: 30)")
	if err := viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout")); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".syncli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".syncli")
	}

	viper.SetEnvPrefix("SYNCLI")

	err := viper.ReadInConfig()

	if err == nil {
		logger.WithFields(logrus.Fields{
			"event": "config_loaded",
			"file":  viper.ConfigFileUsed(),
		}).Info("Configuration file loaded successfully")
	} else {
		logger.WithFields(logrus.Fields{
			"event": "config_load_failed",
			"file":  viper.ConfigFileUsed(),
		}).Warn("Error loading config file, using environment variables")
	}

	viper.AutomaticEnv()

	config.AccessToken = viper.GetString("access_token")
	config.BaseURL = viper.GetString("base_url")

	if config.AccessToken == "" || config.BaseURL == "" {
		logger.WithFields(logrus.Fields{
			"event": "config_validation_failed",
		}).Error("Base URL and Access Token must be provided in the config file or as environment variables")
		os.Exit(1)
	}

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}
}
