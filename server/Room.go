package server

import (
	"chat/service"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Room struct {
	id         string `bson:"id"`
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
	filter := bson.D{{"id", r.id}, {"users.username", client.UserID}}
	err := roomCollection.FindOne(ctx, filter)
	if err == nil {
		r.Clients[client] = true
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	roomCollection := service.GetCollection(service.DB, "rooms")
	filter := bson.D{{"id", r.id}}
	update := bson.D{
		{"$push", bson.D{
			{"sender", message.Sender},
			{"message", message.Message},
		},
		},
	}
	_, err := roomCollection.UpdateOne(ctx, filter, update)
	//TODO - result
	if err == nil {
		for client := range r.Clients {
			client.WriteMess <- *MessageToByte(message)
		}
	} else {
		//TODO - responses
		r.GetClient(message.Sender).WriteMess <- []byte("error")
	}

}
func (r *Room) RunRoom() {
	for {
		select {
		case message := <-r.Broadcast:
			r.sendMessage(message)
		case user := <-r.AddUser:
			r.Clients[user] = true
			mess, _ := json.Marshal(Message{Message: "dołączono do pokoju " + r.id})
			user.WriteMess <- mess
		case key := <-r.RemoveUser:
			delete(r.Clients, key)
		case <-r.Stop:
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}

}
