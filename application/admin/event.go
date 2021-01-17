package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/infrastructure/utils"
	"github.com/sirupsen/logrus"
)

type Event struct{}

func (t *Event) Post(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event post")

	ctx.JSON(status, resp)
}

func (t *Event) Get(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event get")

	ctx.JSON(status, resp)
}

func (t *Event) Put(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event put")

	ctx.JSON(status, resp)
}

func (t *Event) Delete(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event delete")

	ctx.JSON(status, resp)
}

func (t *Event) Status(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event status")

	ctx.JSON(status, resp)
}
