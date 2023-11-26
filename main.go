package main

import (
	"chat/Client"
	"chat/Server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	server := Server.Server{Clients: make(map[*Client.Client]bool)}
	app := gin.Default()
	app.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		defer ws.Close()
		if err != nil {
			log.Println("websocket upgrade error: ", err)
			return
		}
		client := Client.Client{Socket: ws}
		server.AddClient(&client)
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
