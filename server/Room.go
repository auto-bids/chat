package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Room struct {
	id         string
	Clients    map[*Client]bool
	Server     *Server
	Broadcast  chan Message
	AddUser    chan *Client
	RemoveUser chan *Client
}

func (r *Room) AddClient(client *Client) {
	r.Clients[client] = true
}
func (r *Room) RemoveClient(client *Client) {
	delete(r.Clients, client)
}
func CreateRoom(name string, server *Server) *Room {
	return &Room{
		id:         name,
		Server:     server,
		Broadcast:  make(chan Message),
		Clients:    make(map[*Client]bool),
		AddUser:    make(chan *Client),
		RemoveUser: make(chan *Client),
	}
}
func Run(r *Room) {
	for {
		select {
		case message := <-r.Broadcast:
			byteMessage, err := json.Marshal(message)
			for c := range r.Clients {
				if message.Sender != c.UserID {
					if err != nil {
						c.Socket.WriteMessage(websocket.TextMessage, []byte("message error"))
					} else {
						c.Socket.WriteMessage(websocket.TextMessage, byteMessage)
					}
				}

			}
		case user := <-r.AddUser:
			r.Clients[user] = true
			fmt.Println(len(r.Clients))
			fmt.Println("dodano do ", r.id, " uzytkownia: ", user.UserID)
		case key := <-r.RemoveUser:
			delete(r.Clients, key)
		}
	}
}
