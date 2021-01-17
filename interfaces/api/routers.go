package api

import (
	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/application/api"
)

func initRouters(g *gin.Engine) {
	event := g.Group("/event")
	eventApp := api.Event{}
	event.GET("/list", eventApp.List)
	event.GET("/info", eventApp.Info)
	event.POST("/subscribe", eventApp.Subscribe)

	shop := g.Group("/shop")
	shopApp := api.Shop{}
	shop.PUT("/cart/add", shopApp.AddCart)
}
