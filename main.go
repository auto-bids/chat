package main

import (
	"chat/controllers"
	"chat/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	Server := server.CreateServer()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.GET("/chat", func(ctx *gin.Context) {
		controllers.ManageWs(Server, ctx)
	})
	app.POST("/chat/addroom", controllers.CreateRoom)
	app.POST("/chat/adduser", controllers.CreateUser)
	log.Fatal(app.Run(os.Getenv("PORT")))

}
