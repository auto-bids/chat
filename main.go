package main

import (
	"chat/controllers"
	"chat/server"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	Server := server.CreateServer()
	app.GET("/chat", func(ctx *gin.Context) { controllers.ManageWs(Server, ctx) })
	app.GET("/chat/:id", controllers.GetMessages)
	app.POST("/chat/addroom", controllers.CreateRoom)
	app.POST("/chat/adduser", controllers.CreateUser)
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}
