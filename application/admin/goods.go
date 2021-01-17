package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/infrastructure/utils"
	"github.com/sirupsen/logrus"
)

type Goods struct{}

func (t *Goods) Post(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("goods post")

	ctx.JSON(status, resp)
}

func (t *Goods) Get(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("goods get")

	ctx.JSON(status, resp)
}

func (t *Goods) Put(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("goods put")

	ctx.JSON(status, resp)
}

func (t *Goods) Delete(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("goods delete")

	ctx.JSON(status, resp)
}
