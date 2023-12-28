package Client

type Message struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
	Room    string `json:"room"`
	options string `json:"options"`
}
