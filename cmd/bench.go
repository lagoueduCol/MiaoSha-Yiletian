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
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/letian0805/seckill/infrastructure/stores/etcd"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

type Task struct {
	ID          int      `json:"id"`
	Servers     []string `json:"servers"`
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Data        string   `json:"data"`
	ContentType string   `json:"content_type"`
	Concurrency int      `json:"concurrency"`
	Number      int      `json:"number"`
	Duration    int      `json:"duration"`
	Status      int32    `json:"status"`
}

type TaskManager struct {
	sync.Mutex
	task Task
}

func (tm *TaskManager) onConfigChange(task Task) {
	if atomic.LoadInt32(&tm.task.Status) == 1 {
		atomic.StoreInt32(&tm.task.Status, 2)
	}

	for atomic.LoadInt32(&tm.task.Status) == 2 {
		time.Sleep(time.Second)
	}
	tm.Lock()
	tm.task = task
	tm.Unlock()
	if task.Status == 1 {
		tm.task.doBench()
	}
}

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

		// TODO: 还需要加上初始化 etcd 的代码，待完善
		// watchTaskConfig((&TaskManager{}).onConfigChange)
	},
}

var (
	requests    int32
	concurrency int
	url         string
	keepalive   bool
)

func init() {
	rootCmd.AddCommand(benchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	benchCmd.PersistentFlags().Int32VarP(&requests, "requests", "r", 10000, "requests")
	benchCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "C", 50, "concurrency")
	benchCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "url")
	benchCmd.PersistentFlags().BoolVarP(&keepalive, "keepalive", "k", false, "k")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func watchTaskConfig(callback func(cfg Task)) error {
	var err error
	cli := etcd.GetClient()
	key := "/bench/task/config"
	update := func(kv *mvccpb.KeyValue) (bool, error) {
		if string(kv.Key) == key {
			var tmpConfig Task
			err = json.Unmarshal(kv.Value, &tmpConfig)
			if err != nil {
				logrus.Error("update bench config failed, error:", err)
				return false, err
			}
			logrus.Info("update bench config ", tmpConfig)
			callback(tmpConfig)
			return true, nil
		}
		return false, nil
	}
	watchCh := cli.Watch(context.Background(), key)
	for resp := range watchCh {
		for _, evt := range resp.Events {
			if evt.Type == etcdv3.EventTypePut {
				if ok, err := update(evt.Kv); ok {
					break
				} else if err != nil {
					break
				}
			}
		}
	}
	return nil
}

func (t *Task) doBench() {
	wg := &sync.WaitGroup{}
	wg.Add(t.Concurrency)
	for i := 0; i < t.Concurrency; i++ {
		go func() {
			for atomic.LoadInt32(&t.Status) == 1 {
				// do test
			}
			atomic.StoreInt32(&t.Status, 3)
			wg.Done()
		}()
	}
	wg.Wait()
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
					TLSClientConfig: &tls.Config{},
				},
				Timeout: 10 * time.Second,
			}
			wg1.Done()
			<-startCh
			for atomic.AddInt32(&reqs, -1) >= 0 {
				if !keepalive {
					cli = &http.Client{
						Transport: &http.Transport{
							TLSClientConfig:   &tls.Config{},
							DisableKeepAlives: true,
						},
						Timeout: 10 * time.Second,
					}
				}
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
				if !keepalive {
					cli.CloseIdleConnections()
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
