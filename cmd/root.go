/*
Copyright © 2021 letian0805@sina.com

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
	"os"

	"github.com/letian0805/seckill/infrastructure/cluster"
	"github.com/letian0805/seckill/infrastructure/logger"
	"github.com/letian0805/seckill/infrastructure/stores/etcd"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seckill",
	Short: "Seckill server.",
	Long:  `Seckill server.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logrus.Info(cfgFile)
}

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&cfgFile, "config", "c", "./config/seckill.toml", "config file (default is $HOME/.seckill.toml)")
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

		// Search config in home directory with name ".seckill" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".seckill")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
	if err := etcd.Init(); err != nil {
		panic(err)
	}
	logger.InitLogger()
	if err := cluster.WatchClusterConfig(); err != nil {
		panic(err)
	}
}
