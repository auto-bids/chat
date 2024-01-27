package controllers

import (
	"chat/models"
	"chat/responses"
	"chat/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func CreateUser(ctx *gin.Context) {
	result := make(chan responses.Response)
	go func(c *gin.Context) {
		ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer close(result)
		defer cancel()
		var user models.PostUserDB
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := c.ShouldBindJSON(&user); err != nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Data:    map[string]interface{}{"error": err.Error()},
			}
		}
		if err := validate.Struct(user); err != nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Error validation user",
				Data:    map[string]interface{}{"error": err.Error()},
			}
		}
		var userCollection = service.GetCollection(service.DB, "users")
		filter := bson.D{{"email", user.Email}}
		err := userCollection.FindOne(ctxDB, filter)
		if err == nil {
			result <- responses.Response{
				Status:  http.StatusBadRequest,
				Message: "user exists",
				Data:    map[string]interface{}{},
			}
		}
	}(ctx.Copy())

}
