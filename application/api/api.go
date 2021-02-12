package api

import (
	"net/http"
	"time"

	"github.com/letian0805/seckill/domain/user"

	"github.com/letian0805/seckill/infrastructure/utils"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Event struct{}

type Shop struct{}

func (e *Event) List(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	now := time.Now().UnixNano()
	if now%10 == 0 {
		time.Sleep(time.Millisecond * 15)
	}
	//logrus.Info("event list")

	ctx.JSON(status, resp)
}

func (e *Event) Info(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event info")

	ctx.JSON(status, resp)
}

func (e *Event) Subscribe(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("event subscribe")

	ctx.JSON(status, resp)
}

func (s *Shop) AddCart(ctx *gin.Context) {
	resp := &utils.Response{
		Code: 0,
		Data: nil,
		Msg:  "ok",
	}
	status := http.StatusOK

	logrus.Info("shop add cart")

	ctx.JSON(status, resp)
}

type User struct{}

func (u User) Login(ctx *gin.Context) {
	var (
		uid    string
		passwd string
		ok     bool
	)
	if uid, ok = ctx.GetPostForm("uid"); !ok {
		utils.Abort(ctx, http.StatusUnauthorized, "login failed")
		return
	}
	if passwd, ok = ctx.GetPostForm("password"); !ok {
		utils.Abort(ctx, http.StatusUnauthorized, "login failed")
		return
	}
	info, token := user.Login(uid, passwd)
	if info != nil {
		ctx.Header(user.TokenHeader, user.TokenPrefix+token)
		utils.ResponseJSON(ctx, http.StatusOK, "success", nil)
	} else {
		utils.Abort(ctx, http.StatusUnauthorized, "login failed")
	}
}
