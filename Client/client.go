package Client

import "github.com/gorilla/websocket"

type Client struct {
	Socket   *websocket.Conn
	username string
}
