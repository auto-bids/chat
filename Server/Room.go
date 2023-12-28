package Server

import (
	"chat/Client"
	"github.com/gorilla/websocket"
)

type Room struct {
	id         string
	Clients    map[*Client.Client]bool
	Server     *Server
	Broadcast  chan []byte
	AddUser    chan *Client.Client
	RemoveUser chan *Client.Client
}

func (r *Room) AddClient(client *Client.Client) {
	r.Clients[client] = true
}
func (r *Room) RemoveClient(client *Client.Client) {
	delete(r.Clients, client)
}
func CreateRoom(name string, server *Server) *Room {
	return &Room{
		id:         name,
		Server:     server,
		Broadcast:  make(chan []byte),
		Clients:    make(map[*Client.Client]bool),
		AddUser:    make(chan *Client.Client),
		RemoveUser: make(chan *Client.Client),
	}
}
func Run(r *Room) {
	for {
		select {
		case message := <-r.Broadcast:
			for c := range r.Clients {
				c.Socket.WriteMessage(websocket.TextMessage, message)
			}
		case user := <-r.AddUser:
			r.Clients[user] = true
		case key := <-r.RemoveUser:
			delete(r.Clients, key)
		}
	}
}
