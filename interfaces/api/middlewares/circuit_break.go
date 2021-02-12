package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/infrastructure/utils"
)

func NewCircuitBreakMiddleware(cb *utils.CircuitBreaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok := cb.Allow(func() bool {
			c.Next()
			if c.Writer.Status() >= http.StatusInternalServerError {
				return false
			}
			return true
		})
		if !ok {
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	}
}
