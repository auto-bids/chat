package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	Socket     *websocket.Conn
	ReadMess   chan Message
	WrtiteMess chan Message
	Server     *Server
	UserID     string
	Close      chan string
}

func (c *Client) HandleMessages() {
	for {
		_, read, _ := c.Socket.ReadMessage()
		mess := Message{}
		json.Unmarshal(read, &mess)
		fmt.Println(mess)
	}
}

func NewClient(socket *websocket.Conn, ctx *gin.Context) *Client {
	username := ctx.Request.Header["Username"][0]
	fmt.Println(username)
	return &Client{
		Socket: socket,
		Close:  make(chan string),
		UserID: username,
	}
}
func subscribeRoom(roomId string) {

}
