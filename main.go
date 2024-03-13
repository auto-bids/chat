package main

import (
	"chat/routes"
	"chat/server"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	app := gin.Default()
	Server := server.CreateServer()
	routes.ChatRoute(app, Server)
	log.Fatal(app.Run(":" + os.Getenv("PORT")))
}
