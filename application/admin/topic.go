package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/infrastructure/utils"
	"github.com/sirupsen/logrus"
)

type Topic struct{}

func (t *Topic) Post(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("topic post")

	ctx.JSON(status, resp)
}

func (t *Topic) Get(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("topic get")

	ctx.JSON(status, resp)
}

func (t *Topic) Put(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("topic put")

	ctx.JSON(status, resp)
}

func (t *Topic) Delete(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("topic delete")

	ctx.JSON(status, resp)
}

func (t *Topic) Status(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("topic status")

	ctx.JSON(status, resp)
}
