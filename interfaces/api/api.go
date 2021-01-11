package api

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/letian0805/seckill/infrastructure/utils"

	"golang.org/x/sys/unix"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var lis net.Listener

func Run() error {
	var err error

	bind := viper.GetString("api.bind")
	lisCfg := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			logrus.Info("control")
			var err error
			err1 := c.Control(func(fd uintptr) {
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				if err != nil {
					logrus.Error("set socket option failed ", err)
				}
			})
			if err1 != nil {
				logrus.Error("control listener failed ", err1)
				err = err1
			}
			return err
		},
	}
	lis, err = lisCfg.Listen(context.Background(), "tcp", bind)
	if err != nil {
		return err
	}

	g := gin.New()
	// TODO: 初始化路由

	// 更新程序，给老版本发送信号
	go updateProc()

	// 监控黑名单变更
	utils.WatchBlacklist()

	return g.RunListener(lis)
}

func updateProc() {
	if pidFile, err := os.Open(viper.GetString("global.pid")); err == nil {
		pidBytes, _ := ioutil.ReadAll(pidFile)
		pid, _ := strconv.Atoi(string(pidBytes))
		if pid > 0 {
			// 为了避免因某些原因老版本程序无法退出，尝试发送多个信号，最后一次 SIGKILL 将强制结束老版程序
			signals := []syscall.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
			if proc, err := os.FindProcess(pid); err == nil {
				for _, sig := range signals {
					if err = proc.Signal(sig); err != nil {
						break
					}
					var stat *os.ProcessState
					// 等待老版程序退出
					stat, err = proc.Wait()
					if err != nil || stat.Exited() {
						break
					}
				}
			}
		}
		pidFile.Close()
	}
	if pidFile, err := os.Create(viper.GetString("global.pid")); err == nil {
		pid := os.Getpid()
		pidFile.Write([]byte(strconv.Itoa(pid)))
		pidFile.Close()
	}
}

func Exit() {
	lis.Close()
	// TODO: 等待请求处理完
	time.Sleep(10 * time.Second)
}
