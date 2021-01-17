package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/viper"
)

func UpdateProc(service string) {
	fileName := viper.GetString("global.pid")
	fileName = fmt.Sprintf("%s.%s", fileName, service)
	if pidFile, err := os.Open(fileName); err == nil {
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
	if pidFile, err := os.Create(fileName); err == nil {
		pid := os.Getpid()
		pidFile.Write([]byte(strconv.Itoa(pid)))
		pidFile.Close()
	}
}
