package admin

import (
	"net"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/infrastructure/utils"
	"github.com/spf13/viper"
)

var lis net.Listener

func Run() error {
	var err error
	bind := viper.GetString("admin.bind")
	logrus.Info("run admin server on ", bind)
	lis, err = utils.Listen("tcp", bind)
	if err != nil {
		return err
	}

	g := gin.New()

	// 更新程序，给老版本发送信号
	go utils.UpdateProc("admin")

	// 初始化路由
	initRouters(g)
	// 运行服务
	return g.RunListener(lis)
}

func Exit() {
	lis.Close()
	// TODO: 等待请求处理完
	// time.Sleep(10 * time.Second)
	logrus.Info("admin server exit")
}
