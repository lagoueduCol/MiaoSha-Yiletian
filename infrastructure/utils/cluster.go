package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var clusterCfgViper = viper.New()

func WatchClusterConfig() {
	etcdAddr := viper.GetString("etcd.address")
	if etcdAddr == "" {
		return
	}
	v := clusterCfgViper
	if err := v.AddRemoteProvider("etcd", viper.GetString("etcd.address"), "/seckill/config"); err != nil {
		logrus.Error("add remote provider failed, error ", err)
		panic(err)
	}
	go func() {
		if err := v.WatchRemoteConfigOnChannel(); err != nil {
			logrus.Error("watch remote config failed")
			panic(err)
		}
	}()
}
