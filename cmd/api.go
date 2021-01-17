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
	"sync"
	"syscall"

	"github.com/letian0805/seckill/interfaces/rpc"

	"github.com/letian0805/seckill/interfaces/api"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Seckill api server.",
	Long:  `Seckill api server.`,
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		wg.Add(2)
		onApiExit := make(chan error, 1)
		onRpcExit := make(chan error, 1)
		go func() {
			defer wg.Done()
			if err := api.Run(); err != nil {
				logrus.Error(err)
				onApiExit <- err
			}
			close(onApiExit)
		}()
		go func() {
			defer wg.Done()
			if err := rpc.Run(); err != nil {
				logrus.Error(err)
				onRpcExit <- err
			}
			close(onRpcExit)
		}()
		onSignal := make(chan os.Signal)
		signal.Notify(onSignal, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-onSignal:
			logrus.Info("exit by signal ", sig)
			api.Exit()
			rpc.Exit()
		case err := <-onApiExit:
			rpc.Exit()
			logrus.Info("exit by error ", err)
		case err := <-onRpcExit:
			api.Exit()
			logrus.Info("exit by error ", err)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
