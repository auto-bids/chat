package controllers

import (
	"chat/responses"
	"chat/service"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateRoom(c *gin.Context) {
	result := make(chan responses.Response)
	go func(ctx *gin.Context) {
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer close(result)
		var result responses.Response
		//todo
		var collection = service.GetCollection(service.DB, "rooms")
		//todo
	}(c.Copy())
}
