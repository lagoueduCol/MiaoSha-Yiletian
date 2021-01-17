/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"os/signal"
	"syscall"

	"github.com/letian0805/seckill/interfaces/admin"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Seckill admin server.",
	Long:  `Seckill admin server.`,
	Run: func(cmd *cobra.Command, args []string) {
		onExit := make(chan error)
		go func() {
			if err := admin.Run(); err != nil {
				logrus.Error(err)
				onExit <- err
			}
			close(onExit)
		}()
		onSignal := make(chan os.Signal)
		signal.Notify(onSignal, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-onSignal:
			logrus.Info("exit by signal ", sig)
			admin.Exit()
		case err := <-onExit:
			logrus.Info("exit by error ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
}
