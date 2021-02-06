package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Code int         `json:"code"` // 业务错误码
	Data interface{} `json:"data"` // 数据
	Msg  string      `json:"msg"`  // 提示信息
}

func Abort(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(code, &Response{
		Code: code,
		Msg:  msg,
	})
}

func ResponseJSON(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(code, &Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}
