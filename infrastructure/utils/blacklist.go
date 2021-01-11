package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func WatchBlacklist() {
	v := viper.New()
	v.SetConfigFile(viper.GetString("blacklist.filePath"))
	v.OnConfigChange(onBlackListChange)
	go v.WatchConfig()
}

func onBlackListChange(in fsnotify.Event) {
	const writeOrCreateMask = fsnotify.Write | fsnotify.Create
	if in.Op&writeOrCreateMask != 0 {
		updateBlacklist()
	}
}

func updateBlacklist() {
	// TODO: do update
}
