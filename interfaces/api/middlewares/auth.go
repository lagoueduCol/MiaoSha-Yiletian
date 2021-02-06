package middlewares

import (
	"net/http"
	"strings"

	"github.com/letian0805/seckill/infrastructure/utils"

	"github.com/letian0805/seckill/domain/user"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(redirect bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var info *user.Info
		token := ctx.Request.Header.Get(user.TokenHeader)
		if token != "" && strings.Contains(token, user.TokenPrefix) {
			token = strings.Trim(token, user.TokenPrefix)
			token = strings.TrimSpace(token)
			info = user.Auth(token)
		}
		if info != nil {
			ctx.Set("UserInfo", info)
		} else if redirect {
			utils.Abort(ctx, http.StatusUnauthorized, "need login")
			return
		}
		ctx.Next()
	}
}
