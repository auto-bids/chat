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
	"strconv"
	"time"
)

func GetMessages(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		email := ctx.Param("email")
		page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		roomCollection := service.GetCollection(service.DB, "rooms")
		id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		filterUser := bson.D{{"_id", id}, {"users", email}}
		err = roomCollection.FindOne(ctxDB, filterUser).Err()
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		pipeline := bson.A{
			bson.D{{"$match", bson.D{{"_id", id}}}},
			bson.D{{"$unwind", "$messages"}},
			bson.D{{"$sort", bson.D{{"messages.Time", -1}}}},
			bson.D{{"$skip", page * 10}},
			bson.D{{"$limit", 10}},
		}
		cursor, err := roomCollection.Aggregate(ctxDB, pipeline)
		var results []models.MessageUnwindDB
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusNotFound,
				Message: "room not found",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		if err = cursor.All(ctxDB, &results); err != nil {
			panic(err)
		}

		result <- responses.Response{
			Status:  http.StatusFound,
			Message: "room found",
			Data:    map[string]interface{}{"data": results},
		}
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
