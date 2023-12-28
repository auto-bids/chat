package Server

type Server struct {
	Clients map[string]*Room
}

func CreateServer() *Server {
	return &Server{
		Clients: make(map[string]*Room),
	}
}
func (s Server) AddRoom(name string) {
	room := CreateRoom(name, &s)
	s.Clients[name] = room
	go Run(room)
}
