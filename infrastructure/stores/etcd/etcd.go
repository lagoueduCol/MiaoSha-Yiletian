package etcd

import (
	"sync"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

var etcdCli *etcd.Client
var etcdOnce = &sync.Once{}

func Init() error {
	var err error
	etcdOnce.Do(
		func() {
			endpoints := viper.GetStringSlice("etcd.endpoints")
			username := viper.GetString("etcd.username")
			password := viper.GetString("etcd.password")
			cfg := etcd.Config{
				Endpoints:   endpoints,
				DialTimeout: time.Second,
				Username:    username,
				Password:    password,
			}
			etcdCli, err = etcd.New(cfg)
		})
	return err
}

func GetClient() *etcd.Client {
	return etcdCli
}
