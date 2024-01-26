package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	Socket    *websocket.Conn
	WriteMess chan []byte
	Server    *Server
	UserID    string
	Rooms     map[string]*Room
	Close     chan string
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewClient(socket *websocket.Conn, ctx *gin.Context) *Client {
	username := ctx.Request.Header["Username"][0]
	return &Client{
		Socket:    socket,
		Close:     make(chan string),
		Rooms:     make(map[string]*Room),
		WriteMess: make(chan []byte),
		UserID:    username,
	}
}
func (c *Client) sendToServer(roomId string, mess *Message) {
	c.Rooms[roomId].Broadcast <- mess
}
func (c *Client) subscribeRoom(roomId string) error {
	room := c.Server.GetRoom(roomId)
	if room == nil {
		room = c.Server.AddRoom(roomId)
	}
	c.Rooms[roomId] = room
	room.AddUser <- c
	return nil
}
func (c *Client) unsubscribeRoom(roomId string) error {
	room := c.Server.GetRoom(roomId)
	delete(c.Rooms, roomId)
	room.RemoveUser <- c
	return nil
}

func (c *Client) ReadPump() {
	defer func() {
		c.Socket.Close()
	}()
	c.Socket.SetReadLimit(maxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error { c.Socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		mess := ByteToMessage(message)
		switch mess.Options {
		case "subscribe":
			c.subscribeRoom(mess.Destination)

		case "unsubscribe":
			c.unsubscribeRoom(mess.Destination)
		case "message":
			c.sendToServer(mess.Destination, mess)
		default:
		}
		time.Sleep(time.Millisecond)
	}
}
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.WriteMess:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.WriteMess)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.WriteMess)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
		time.Sleep(time.Millisecond)
	}
}
