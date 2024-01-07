package server

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
func (s *Server) AddRoom(roomId string) *Room {
	room := CreateRoom(roomId, s)
	go room.RunRoom()
	s.Rooms[room.id] = room
	return room
}
func (s *Server) RemoveRoom(roomId string) {
	s.Rooms[roomId].Stop <- true
	delete(s.Rooms, roomId)
}
func (s *Server) GetRoom(name string) *Room {
	var room *Room
	room = s.Rooms[name]
	if room == nil {
		return nil
	}
	return room
}
func (s *Server) AddClient(client *Client) {
	s.Client[client] = true
}
func (s *Server) AddClientToRoom(r *Room, c *Client) {
	r.AddUser <- c
}
