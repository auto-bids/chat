package Server

type Server struct {
	Rooms      map[string]*Room
	RemoveRoom chan string
	AddRoom    chan string
}
