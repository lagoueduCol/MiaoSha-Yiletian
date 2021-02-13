package api

import (
	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/application/api"
	"github.com/letian0805/seckill/infrastructure/utils"
	"github.com/letian0805/seckill/interfaces/api/middlewares"
)

func initRouters(g *gin.Engine) {
	g.POST("/login", api.User{}.Login)

	eventCB := utils.NewCircuitBreaker(
		utils.WithDuration(100),
		utils.WithTotalLimit(2000),
		utils.WithLatencyLimit(100),
		utils.WithFailsLimit(5),
	)
	eventCBMdw := middlewares.NewCircuitBreakMiddleware(eventCB)
	event := g.Group("/event").Use(eventCBMdw, middlewares.NewAuthMiddleware(false))
	eventApp := api.Event{}
	event.GET("/list", eventApp.List)
	event.GET("/info", eventApp.Info)

	subscribe := g.Group("/event/subscribe").Use(middlewares.NewAuthMiddleware(true))
	subscribe.POST("/", eventApp.Subscribe)
	//
	//shopCB := utils.NewCircuitBreaker(
	//	utils.WithDuration(100),
	//	utils.WithTotalLimit(1000),
	//	utils.WithLatencyLimit(200),
	//	utils.WithFailsLimit(5),
	//)
	//shopCBMdw := middlewares.NewCircuitBreakMiddleware(shopCB)
	shop := g.Group("/shop") //.Use(shopCBMdw, middlewares.NewAuthMiddleware(true), middlewares.Blacklist)
	shopApp := api.Shop{}
	shop.PUT("/cart/add", shopApp.AddCart)
}
