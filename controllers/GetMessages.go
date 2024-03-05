package controllers

import (
	"chat/models"
	"chat/responses"
	"chat/service"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetMessages(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {

		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		var room models.RoomDB
		roomCollection := service.GetCollection(service.DB, "rooms")
		id, err := primitive.ObjectIDFromHex(ctx.Query("id"))
		email := ctx.Query("email")

		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
		}
		filter := bson.D{{"_id", id}, {"users", bson.D{{"$elemMatch", email}}}}
		err = roomCollection.FindOne(ctxDB, filter).Decode(&room)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusNotFound,
				Message: "room not found",
				Data:    map[string]interface{}{"error": err.Error()},
			}
		}
		result <- responses.Response{
			Status:  http.StatusFound,
			Message: "room found",
			Data:    map[string]interface{}{"data": room.Messages},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
