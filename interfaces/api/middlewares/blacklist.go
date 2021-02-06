package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/domain/user"
	"github.com/letian0805/seckill/infrastructure/utils"
)

func Blacklist(ctx *gin.Context) {
	data, _ := ctx.Get("UserInfo")
	info, ok := data.(*user.Info)
	if !ok {
		utils.Abort(ctx, http.StatusUnauthorized, "need login")
		return
	}
	if utils.InBlacklist(info.UID) {
		utils.Abort(ctx, http.StatusForbidden, "blocked")
		return
	}
	ctx.Next()
}
