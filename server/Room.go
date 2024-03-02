package server

import (
	"chat/models"
	"chat/service"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Room struct {
	id         string
	Clients    map[*Client]bool
	Server     *Server
	Stop       chan bool
	Broadcast  chan *Message
	AddUser    chan *Client
	RemoveUser chan *Client
}

func CreateRoom(name string, server *Server) *Room {
	return &Room{
		id:         name,
		Server:     server,
		Broadcast:  make(chan *Message),
		Clients:    make(map[*Client]bool),
		AddUser:    make(chan *Client),
		RemoveUser: make(chan *Client),
		Stop:       make(chan bool),
	}
}

func (r *Room) AddClient(client *Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	roomCollection := service.GetCollection(service.DB, "rooms")
	id, _ := primitive.ObjectIDFromHex(r.id)
	filter := bson.D{{"_id", id}, {"users", client.UserID}}
	var room models.RoomDB
	err := roomCollection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		client.WriteMess <- []byte("unauthorized")
		return
	}
	client.WriteMess <- []byte(r.id)
	r.Clients[client] = true
}
func (r *Room) RemoveClient(client *Client) {
	delete(r.Clients, client)
}
func (r *Room) GetClient(client string) *Client {
	for i := range r.Clients {
		if i.UserID == client {
			return i
		}
	}
	return nil
}
func (r *Room) sendMessage(message *Message) {
	if r.GetClient(message.Sender) != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		roomCollection := service.GetCollection(service.DB, "rooms")
		id, _ := primitive.ObjectIDFromHex(r.id)
		filter := bson.D{{"_id", id}}
		var room models.RoomDB
		roomCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&room)
		data := models.MessageDB{Sender: message.Sender, Message: message.Message, Time: time.Now().Unix()}
		update := bson.M{"$push": bson.M{"messages": data}}
		_, err := roomCollection.UpdateOne(ctx, filter, update)
		if err == nil {
			for client := range r.Clients {
				client.WriteMess <- *MessageToByte(message)
			}
		} else {
			r.GetClient(message.Sender).WriteMess <- []byte("error")
		}
	}
}
func (r *Room) RunRoom() {
	for {
		select {
		case message := <-r.Broadcast:
			r.sendMessage(message)
		case user := <-r.AddUser:
			r.AddClient(user)
		case key := <-r.RemoveUser:
			delete(r.Clients, key)
			if len(r.Clients) == 0 {
				return
			}
		case <-r.Stop:
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}

}
