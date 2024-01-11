package controllers

import (
	"chat/server"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ManageWs(s *server.Server, ctx *gin.Context, db *mongo.Client) {
	ws, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := server.NewClient(ws, ctx)
	s.AddClient(client)
	client.Server = s
	go client.WritePump()
	go client.ReadPump()
}
