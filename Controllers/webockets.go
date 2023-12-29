package Controllers

import (
	"chat/Server"
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

func ManageWs(s *Server.Server, ctx *gin.Context) {
	room := ctx.Request.Header["Room"][0]
	ws, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("websocket upgrade error: ", err)
		return
	}
	client := Server.NewClient(ws, ctx)
	go client.HandleMessages()
	s.AddClientToRoom(s.GetRoom(room), client)
}
