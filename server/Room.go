package server

import (
	"encoding/json"
	"fmt"
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
	r.Clients[client] = true
}
func (r *Room) RemoveClient(client *Client) {
	delete(r.Clients, client)
}
func (r *Room) RunRoom() {
	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println(message.Message)
			for client := range r.Clients {
				client.WriteMess <- *MessageToByte(message)
			}
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
