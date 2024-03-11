package controllers

import (
	"chat/models"
	"chat/responses"
	"chat/service"
	"context"
	"fmt"
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
		//var room models.RoomDB
		roomCollection := service.GetCollection(service.DB, "rooms")
		id, err := primitive.ObjectIDFromHex(ctx.Query("id"))

		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
		}
		pipeline := bson.A{
			bson.D{{"$match", bson.D{{"_id", id}}}},
			bson.D{{"$unwind", "$messages"}},
			bson.D{{"$sort", bson.D{{"messages.Time", 1}}}},
			bson.D{{"$skip", page * 10}},
			bson.D{{"$limit", 10}},
		}
		cursor, err := roomCollection.Aggregate(ctxDB, pipeline)
		var results []models.MessageUnwindDB
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				panic(err)
			}
			fmt.Println(result)
		}
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusNotFound,
				Message: "room not found",
				Data:    map[string]interface{}{"error": err.Error()},
			}
		}
		if err = cursor.All(context.TODO(), &results); err != nil {
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
