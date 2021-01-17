package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/letian0805/seckill/application/admin"
)

func initRouters(g *gin.Engine) {
	topic := g.Group("/topic")
	topicApp := admin.Topic{}
	topic.POST("/", topicApp.Post)
	topic.GET("/", topicApp.Get)
	topic.GET("/:id", topicApp.Get)
	topic.PUT("/:id", topicApp.Put)
	topic.PUT("/:id/:status", topicApp.Status)
	topic.DELETE("/:id", topicApp.Delete)

	event := g.Group("/event")
	eventApp := admin.Event{}
	event.POST("/", eventApp.Post)
	event.GET("/", eventApp.Get)
	event.GET("/:id", eventApp.Get)
	event.PUT("/:id", eventApp.Put)
	event.PUT("/:id/:status", eventApp.Status)
	event.DELETE("/:id", eventApp.Delete)

	goods := g.Group("/goods")
	goodsApp := admin.Goods{}
	goods.POST("/", goodsApp.Post)
	goods.GET("/", goodsApp.Get)
	goods.GET("/:id", goodsApp.Get)
	goods.PUT("/:id", goodsApp.Put)
	goods.DELETE("/:id", goodsApp.Delete)
}
