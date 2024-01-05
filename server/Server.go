package server

import (
	"errors"
)

type Server struct {
	Rooms  map[string]*Room
	Client map[*Client]bool
}

func CreateServer() *Server {
	return &Server{
		Rooms:  make(map[string]*Room),
		Client: make(map[*Client]bool),
	}
}
func (s *Server) AddRoom(room *Room) {
	s.Rooms[room.id] = room
}
func (s *Server) GetRoom(name string) (*Room, error) {
	var room *Room
	room = s.Rooms[name]
	if room == nil {
		return nil, errors.New("no such room")
	}
	return room, nil
}
func (s *Server) AddClient(client *Client) {
	s.Client[client] = true
}
func (s *Server) AddClientToRoom(r *Room, c *Client) {
	r.AddUser <- c
}
