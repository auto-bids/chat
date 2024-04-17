package server

import (
	"chat/models"
	"chat/service"
	"context"
	"fmt"
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
		UserID:    ctx.Param("email"),
	}
}
func (c *Client) sendToRoom(mess *Message) {
	if c.Rooms[mess.Destination] == nil {
		c.WriteMess <- []byte("unauthorized")
		return
	}
	c.Rooms[mess.Destination].Broadcast <- mess
}
func (c *Client) createRoom(target string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	roomCollection := service.GetCollection(service.DB, "rooms")
	var usersCollection = service.GetCollection(service.DB, "users")
	var user models.PostUserDB
	errFind := usersCollection.FindOne(ctx, bson.D{{"email", target}}).Decode(&user)
	if errFind != nil {
		var user models.PostUserDB
		user.Email = target
		user.Rooms = []models.UserRooms{}
		_, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			c.WriteMess <- []byte(err.Error())
		}
	}
	filter := bson.D{{"users", bson.D{{"$all", bson.A{c.UserID, target}}}}}
	var room models.RoomDB
	err := roomCollection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		var roomInsert models.PostRoomDB
		roomInsert.Users = append(room.Users, c.UserID, target)
		roomInsert.Messages = []models.MessageDB{}
		res, errIns := roomCollection.InsertOne(ctx, roomInsert)
		if errIns != nil {
			return errIns
		}
		id := res.InsertedID.(primitive.ObjectID)
		if target != c.UserID {
			userRoom := models.UserRooms{Id: id.Hex(), Email: target}
			update := bson.M{"$push": bson.M{"rooms": userRoom}}
			_, err = usersCollection.UpdateOne(ctx, bson.D{{"email", c.UserID}}, update)
			fmt.Println(userRoom, update)
			if err != nil {
				c.WriteMess <- []byte(err.Error())
			}
			userRoom = models.UserRooms{Id: id.Hex(), Email: c.UserID}
			update = bson.M{"$push": bson.M{"rooms": userRoom}}
			_, err = usersCollection.UpdateOne(ctx, bson.D{{"email", target}}, update)
			if err != nil {
				c.WriteMess <- []byte(err.Error())
			}
		}
		room.Id = id.Hex()
	}
	roomDB := c.Server.GetRoom(room.Id)
	if roomDB == nil {
		c.Server.AddRoom(room.Id)
		roomDB = c.Server.GetRoom(room.Id)
	}
	c.Rooms[room.Id] = roomDB
	roomDB.AddUser <- c
	return nil
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
		case "create":
			c.createRoom(mess.Destination)
		}
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
	}
}
