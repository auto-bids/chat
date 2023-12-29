package Server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	Socket *websocket.Conn
	Room   *Room
	Close  chan string
	userID string
}

func (c *Client) HandleMessages() {
	for {
		_, read, _ := c.Socket.ReadMessage()
		c.Room.Broadcast <- read
	}
}

func NewClient(socket *websocket.Conn, ctx *gin.Context) *Client {
	username := ctx.Request.Header["Username"][0]
	fmt.Println(username)
	return &Client{
		Socket: socket,
		Close:  make(chan string),
		userID: username,
	}

}
