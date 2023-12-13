package main

import (
	"chat/Client"
	"chat/Server"
	"chat/Controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	server := Server.Server{Clients: make(map[*Client.Client]bool)}
	app := gin.Default()
	app.GET("/ws/:id", Controllers.JoinRoom {

		server.AddClient(&Client.Client{Socket: ws})
		for {
			mt, message, err := ws.ReadMessage()
			server.Broadcast(mt, message)
			if err != nil {
				log.Println("websocket message reading error: ", err)
				break
			}
		}
	})
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Fatal(app.Run(os.Getenv("PORT")))
}
