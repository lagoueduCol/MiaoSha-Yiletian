package redis

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

const Nil = redis.Nil

var cli *redis.Client

func Init() error {
	addr := viper.GetString("redis.address")
	auth := viper.GetString("redis.auth")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	opt := &redis.Options{
		Network:  "tcp",
		Addr:     addr,
		Password: auth,
	}
	cli = redis.NewClient(opt)
	if cli == nil {
		return errors.New("init redis client failed")
	}
	return nil
}

func GetClient() *redis.Client {
	return cli
}
