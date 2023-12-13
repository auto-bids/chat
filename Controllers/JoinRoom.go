package Controllers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func JoinRoom(ctx *gin.Context) {
	go func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		defer ws.Close()
		if err != nil {
			log.Println("websocket upgrade error: ", err)
			return
		}
	}(ctx.Copy())
}

{

server.AddClient(&Client.Client{Socket: ws})
for {
mt, message, err := ws.ReadMessage()
server.Broadcast(mt, message)
if err != nil {
log.Println("websocket message reading error: ", err)
break
		}
	}
}