package main

import (
	"chat/Controllers"
	"chat/Server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	server := Server.CreateServer()
	server.AddRoom("test1")
	server.AddRoom("test2")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.GET("/chat", func(ctx *gin.Context) {
		ws, err := Controllers.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		defer ws.Close()
		if err != nil {
			log.Println("websocket upgrade error: ", err)
			return
		}
	})

	log.Fatal(app.Run(os.Getenv("PORT")))
}
