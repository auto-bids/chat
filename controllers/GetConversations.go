package controllers

import (
	"chat/models"
	"chat/responses"
	"chat/service"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func GetConversations(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		email := ctx.Param("email")
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		userCollection := service.GetCollection(service.DB, "users")
		var user models.UserDB
		filter := bson.D{{"email", email}}
		err := userCollection.FindOne(ctxDB, filter).Decode(&user)
		if err != nil {
			result <- responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Invalid Id",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		result <- responses.Response{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    map[string]interface{}{"data": user},
		}
		return
	}(ctx.Copy())
	res := <-result
	ctx.JSON(res.Status, res)
}
