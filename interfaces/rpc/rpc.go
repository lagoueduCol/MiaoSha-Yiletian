package rpc

import (
	"github.com/letian0805/seckill/infrastructure/utils"

	"github.com/letian0805/seckill/application/api"
	"github.com/letian0805/seckill/application/api/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var grpcS *grpc.Server

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
	return grpcS.Serve(lis)
}

func Exit() {
	grpcS.GracefulStop()
	logrus.Info("rpc server exit")
}
