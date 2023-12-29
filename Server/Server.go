package Server

type Server struct {
	Clients map[string]*Room
}

func CreateServer() *Server {
	return &Server{
		Clients: make(map[string]*Room),
	}
}
func (s *Server) AddRoom(name string) {
	room := CreateRoom(name, s)
	s.Clients[name] = room
	go Run(room)
}
func (s *Server) GetRoom(name string) *Room {
	var Room *Room
	Room = s.Clients[name]
	return Room
}
func (s *Server) AddClientToRoom(r *Room, c *Client) {
	c.Room = r
	r.AddUser <- c
}
