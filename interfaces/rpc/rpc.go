package rpc

import (
	"sync"

	"github.com/letian0805/seckill/infrastructure/cluster"
	"github.com/letian0805/seckill/infrastructure/utils"

	"github.com/letian0805/seckill/application/api"
	"github.com/letian0805/seckill/application/api/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	grpcS *grpc.Server
	once  = &sync.Once{}
	node  *cluster.Node
)

func Run() error {
	bind := viper.GetString("api.rpc")
	logrus.Info("run RPC server on ", bind)
	lis, err := utils.Listen("tcp", bind)
	if err != nil {
		return err
	}

	grpcS = grpc.NewServer()
	eventRPC := &api.EventRPCServer{}
	rpc.RegisterEventRPCServer(grpcS, eventRPC)
	// 支持 gRPC reflection，方便调试
	reflection.Register(grpcS)

	//初始化集群
	cluster.Init("seckill")
	var addr string
	if addr, err = utils.Extract(bind); err == nil {
		//注册节点信息
		version := viper.GetString("api.version")
		if version == "" {
			version = "v0.1"
		}
		once.Do(func() {
			node = &cluster.Node{
				Addr:    addr,
				Version: version,
				Proto:   "gRPC",
			}
			err = cluster.Register(node, 6)
		})
	}

	if err != nil {
		return err
	}

	return grpcS.Serve(lis)
}

func Exit() {
	cluster.Deregister(node)
	grpcS.GracefulStop()
	logrus.Info("rpc server exit")
}
