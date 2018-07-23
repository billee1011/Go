// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.a
// Unless requipache.org/licenses/LICENSE-2.0
////red by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	//"steve/serviceloader/logger"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"steve/servicelauncher/launcher"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "servicelauncher",
	Short: "Debug service",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		launcher.Init(args)
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
	cobra.OnInitialize(_init)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.serviceloader.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func _init() {
	configFile := initConfig()
	initDefaultConfig()
	initLogger()

	if configFile != "" {
		logrus.WithField("config", configFile).Info("using config file")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() string {
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

		// Search config in home directory with name ".serviceloader" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".serviceloader")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		return viper.ConfigFileUsed()
	}
	return ""
}

func initDefaultConfig() {
	viper.SetDefault("log_level", "debug")
	viper.SetDefault("log_dir", "")
	viper.SetDefault("log_prefix", "")
	viper.SetDefault("log_stderr", true)
	viper.SetDefault("rpc_certi_file", "")
	viper.SetDefault("rpc_key_file", "")
	viper.SetDefault("rpc_addr", "")
	viper.SetDefault("rpc_port", 0)
	viper.SetDefault("rpc_ca_file", "")
	viper.SetDefault("rpc_server_name", "")
	viper.SetDefault("certi_server_name", "")
	viper.SetDefault("redis_addr", "127.0.0.1:6379")
	viper.SetDefault("redis_passwd", "")
	viper.SetDefault("consul_addr", "127.0.0.1:8500")
}

func initLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	//logger.SetupLog(viper.GetString("log_prefix"), viper.GetString("log_dir"),
	//	viper.GetString("log_level"), viper.GetBool("log_stderr"))
}
