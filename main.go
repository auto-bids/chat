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
	app.GET("/chat/:username/:email", func(ctx *gin.Context) { controllers.ManageWs(Server, ctx) })
	app.GET("/chat/messages/:id/:email", controllers.GetMessages)
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}
