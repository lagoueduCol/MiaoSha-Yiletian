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
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// benchCmd represents the bench command
var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Seckill benchmark tool.",
	Long:  `Seckill benchmark tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bench called")
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Info("requests ", requests, " concurrency ", concurrency, " url ", url)
		doBench()
	},
}

var requests int32
var concurrency int
var url string

func init() {
	rootCmd.AddCommand(benchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	benchCmd.PersistentFlags().Int32VarP(&requests, "requests", "r", 10000, "requests")
	benchCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "C", 50, "concurrency")
	benchCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "url")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func doBench() {
	runtime.GOMAXPROCS(4)
	wg := &sync.WaitGroup{}
	startCh := make(chan struct{})
	success := int32(0)
	failed := int32(0)
	reqs := requests
	wg.Add(concurrency)
	wg1 := &sync.WaitGroup{}
	wg1.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			cli := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: false,
					},
				},
				Timeout: 10 * time.Second,
			}
			wg1.Done()
			<-startCh
			for atomic.AddInt32(&reqs, -1) >= 0 {
				resp, err := cli.Get(url)
				if err != nil || resp.StatusCode > 404 {
					logrus.Error(err)
					atomic.AddInt32(&failed, 1)
				} else {
					atomic.AddInt32(&success, 1)
				}
				if resp != nil {
					ioutil.ReadAll(resp.Body)
					resp.Body.Close()
				}
			}
			wg.Done()
		}()
	}
	wg1.Wait()
	close(startCh)
	start := time.Now().Unix()
	wg.Wait()
	end := time.Now().Unix()
	fmt.Printf("total: %d, cost: %d, success: %d, failed: %d, qps: %d\n", requests, end-start, success, failed, requests/int32(end-start))
}
