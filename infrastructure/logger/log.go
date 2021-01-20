package logger

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitLogger() {
	logLevel := strings.ToLower(viper.GetString("global.logLevel"))
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else if logLevel == "info" {
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.ErrorLevel)
	}
	formatter := &logrus.JSONFormatter{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyFile:  "file",
			logrus.FieldKeyLevel: "level",
		},
		PrettyPrint: false,
	}

	logrus.SetFormatter(formatter)
}
