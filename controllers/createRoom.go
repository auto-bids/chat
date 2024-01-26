package controllers

import (
	"chat/models"
	"chat/responses"
	"chat/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateRoom(c *gin.Context) {
	result := make(chan responses.Response)
	go func(ctx *gin.Context) {
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer close(result)
		var room models.PostRoomDB
		if err := ctx.ShouldBindJSON(&room); err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "invalid request body",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		var collection = service.GetCollection(service.DB, "rooms")
		res, err := collection.InsertOne(ctxDB, room)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error adding room",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		result <- responses.Response{
			Status:  http.StatusCreated,
			Message: "room created",
			Data:    map[string]interface{}{"data": res},
		}

	}(c.Copy())
	res := <-result
	c.JSON(res.Status, res)
}
