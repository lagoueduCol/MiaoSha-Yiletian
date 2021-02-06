package api

import (
	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/application/api"
	"github.com/letian0805/seckill/interfaces/api/middlewares"
)

func initRouters(g *gin.Engine) {
	g.POST("/login", api.User{}.Login)

	event := g.Group("/event").Use(middlewares.NewAuthMiddleware(false))
	eventApp := api.Event{}
	event.GET("/list", eventApp.List)
	event.GET("/info", eventApp.Info)

	subscribe := g.Group("/event/subscribe").Use(middlewares.NewAuthMiddleware(true))
	subscribe.POST("/", eventApp.Subscribe)

	shop := g.Group("/shop").Use(middlewares.NewAuthMiddleware(true), middlewares.Blacklist)
	shopApp := api.Shop{}
	shop.PUT("/cart/add", shopApp.AddCart)
}
