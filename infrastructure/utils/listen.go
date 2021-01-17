package utils

import (
	"context"
	"net"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func Listen(network string, addr string) (net.Listener, error) {
	lisCfg := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
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
	return lisCfg.Listen(context.Background(), "tcp", addr)
}
