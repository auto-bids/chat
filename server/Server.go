package server

type Server struct {
	Rooms  map[string]*Room
	Client map[string]*Client
}

func CreateServer() *Server {
	return &Server{
		Rooms: make(map[string]*Room),
	}
}
func (s *Server) AddRoom(name string) {
	room := CreateRoom(name, s)
	s.Rooms[name] = room
	go Run(room)
}
func (s *Server) GetRoom(name string) *Room {
	var Room *Room
	Room = s.Rooms[name]
	return Room
}
func (s *Server) AddClient(client *Client) {
	s.Client[client.UserID] = client
}
func (s *Server) AddClientToRoom(r *Room, c *Client) {
	r.AddUser <- c
}
