package Server

import (
	"chat/Client"
)

type Room struct {
	id         string
	Broadcast  chan []byte
	Clients    map[*Client.Client]bool
	AddUser    chan *Client.Client
	RemoveUser chan *Client.Client
}

func (r *Room) AddClient(client *Client.Client) {
	r.Clients[client] = true
}
func (r *Room) RemoveClient(client *Client.Client) {
	delete(r.Clients, client)
}
func run(r *Room) {
	for {
		select {
		case message := <-r.Broadcast:
			for c := range r.Clients {
				c.Socket.WriteMessage(message)
			}
		}
	}
}
