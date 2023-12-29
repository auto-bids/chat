package Server

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Room struct {
	id         string
	Clients    map[*Client]bool
	Server     *Server
	Broadcast  chan []byte
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
		Broadcast:  make(chan []byte),
		Clients:    make(map[*Client]bool),
		AddUser:    make(chan *Client),
		RemoveUser: make(chan *Client),
	}
}
func Run(r *Room) {
	for {
		select {
		case message := <-r.Broadcast:
			fmt.Println(string(message))
			for c := range r.Clients {
				c.Socket.WriteMessage(websocket.TextMessage, message)
			}
		case user := <-r.AddUser:
			r.Clients[user] = true
			fmt.Println(len(r.Clients))
			fmt.Println("dodano do ", r.id, " uzytkownia: ", user.userID)
		case key := <-r.RemoveUser:
			delete(r.Clients, key)
		}
	}
}
