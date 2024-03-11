package server

import "sync"

type Server struct {
	Mutex   sync.Mutex
	Rooms   map[string]*Room
	Clients map[*Client]bool
}

func CreateServer() *Server {
	return &Server{
		Rooms:   make(map[string]*Room),
		Clients: make(map[*Client]bool),
	}
}
func (s *Server) AddRoom(id string) *Room {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	room := CreateRoom(id, s)
	go room.RunRoom()
	s.Rooms[id] = room
	return room
}
func (s *Server) RemoveRoom(id string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Rooms[id].Stop <- true
	delete(s.Rooms, id)
}
func (s *Server) GetRoom(id string) *Room {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	var room *Room
	room = s.Rooms[id]
	if room == nil {
		return nil
	}
	return room
}
func (s *Server) AddClient(client *Client) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Clients[client] = true
}
