package controllers

import (
	"chat/server"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ManageWs(s *server.Server, ctx *gin.Context) {
	ws, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("websocket upgrade error: ", err)
		return
	}
	client := server.NewClient(ws, ctx)
	go client.HandleMessages()
}
