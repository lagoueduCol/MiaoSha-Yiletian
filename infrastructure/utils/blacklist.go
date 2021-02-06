package utils

import (
	"bufio"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var blacklist struct {
	sync.RWMutex
	data map[string]struct{}
}

func init() {
	blacklist.data = make(map[string]struct{})
}

func WatchBlacklist() {
	v := viper.New()
	v.SetConfigFile(viper.GetString("blacklist.filePath"))
	v.OnConfigChange(onBlacklistChange)
	go v.WatchConfig()
}

func onBlacklistChange(in fsnotify.Event) {
	const writeOrCreateMask = fsnotify.Write | fsnotify.Create
	if in.Op&writeOrCreateMask != 0 {
		updateBlacklist()
	}
}

func updateBlacklist() {
	filePath := viper.GetString("blacklist.filePath")
	fp, err := os.Open(filePath)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer fp.Close()

	data := make(map[string]struct{})
	f := bufio.NewReader(fp)
	for {
		line, _, err := f.ReadLine()
		if err != nil {
			break
		}
		data[string(line)] = struct{}{}
	}
	blacklist.Lock()
	blacklist.data = data
	blacklist.Unlock()
}

func InBlacklist(uid string) bool {
	blacklist.RLock()
	_, ok := blacklist.data[uid]
	blacklist.RUnlock()
	return ok
}
