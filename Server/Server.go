package Server

import (
	"chat/Client"
	"log"
)

type Server struct {
	Clients map[*Client.Client]bool
}

func (s *Server) Broadcast(messageType int, message []byte) {

	for client := range s.Clients {
		err := client.Socket.WriteMessage(messageType, message)
		if err != nil {
			log.Println("websocket broadcasting error: ", err)
		}
	}
}
func (s *Server) AddClient(client *Client.Client) {
	s.Clients[client] = true
}
func (s *Server) RemoveClient(client *Client.Client) {
	delete(s.Clients, client)
}
