package Client

import (
	"chat/Server"
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	Socket *websocket.Conn
	Room   *Server.Room
	Close  chan string
	userID string
}

func writeMess(c *Client) {

}
func handleMessages(c *Client) {
	for {
		select {
		case read := <-c.Room.Broadcast:
			fmt.Println(string(read))
		}

	}
}
