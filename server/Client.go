package server

import (
	"chat/models"
	"chat/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	return &Client{
		Socket:    socket,
		Close:     make(chan string),
		Rooms:     make(map[string]*Room),
		WriteMess: make(chan []byte),
		UserID:    ctx.Query("email"),
	}
}
func (c *Client) sendToRoom(mess *Message) {
	if c.Rooms[mess.Destination] == nil {
		c.WriteMess <- []byte("unauthorized")
		return
	}
	c.Rooms[mess.Destination].Broadcast <- mess
}
func (c *Client) subscribeRoom(roomId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var user models.UserDB
	roomCollection := service.GetCollection(service.DB, "rooms")
	id, _ := primitive.ObjectIDFromHex(roomId)
	filter := bson.D{{"_id", id}, {"users", c.UserID}}
	err := roomCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}
	room := c.Server.GetRoom(roomId)
	if room == nil {
		c.Server.AddRoom(roomId)
		room = c.Server.GetRoom(roomId)
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
func (c *Client) closeConnection() {
	for _, room := range c.Rooms {
		room.RemoveUser <- c
	}
	delete(c.Server.Clients, c)
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
				c.closeConnection()
			}
			break
		}
		mess := ByteToMessage(message)
		mess.Sender = c.UserID
		switch mess.Options {
		case "subscribe":
			c.subscribeRoom(mess.Destination)
		case "unsubscribe":
			c.unsubscribeRoom(mess.Destination)
		case "message":
			c.sendToRoom(mess)
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
