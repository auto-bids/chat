package routes

import (
	"chat/controllers"
	"chat/server"
	"github.com/gin-gonic/gin"
)

func ChatRoute(router *gin.Engine, Server *server.Server) {
	chat := router.Group("/chat")
	{
		chat.GET("/:email", func(ctx *gin.Context) { controllers.ManageWs(Server, ctx) })
		chat.GET("/messages/:id/:email/:page", controllers.GetMessages)
		chat.GET("/conversations/:email", controllers.GetConversations)
	}
}
