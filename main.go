package main

import (
	"chat/controllers"
	"chat/server"
	"chat/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	Server := server.CreateServer()
	db := service.ConnectDB()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.GET("/chat", func(ctx *gin.Context) {
		controllers.ManageWs(Server, ctx, db)
	})

	log.Fatal(app.Run(os.Getenv("PORT")))
}
