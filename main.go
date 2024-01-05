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
	Ra := server.CreateRoom("test1", Server)
	Rb := server.CreateRoom("test2", Server)
	go Rb.RunRoom()
	go Ra.RunRoom()
	Server.AddRoom(Ra)
	Server.AddRoom(Rb)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.GET("/chat", func(ctx *gin.Context) {
		controllers.ManageWs(Server, ctx)
	})
	log.Fatal(app.Run(os.Getenv("PORT")))
}
