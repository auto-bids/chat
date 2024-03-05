package controllers

import (
	"chat/models"
	"chat/responses"
	"chat/server"
	"chat/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func addUserToDb(ctx *gin.Context) (*mongo.InsertOneResult, error) {
	email := ctx.Query("email")
	username := ctx.Query("username")
	fmt.Println(email, username)
	ctxDB, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var user models.PostUserDB
	user.Email = email
	user.Username = username
	user.Rooms = []string{}
	var userCollection = service.GetCollection(service.DB, "users")
	filter := bson.D{{"email", user.Email}, {"username", user.Username}}
	var existingUser models.UserDB
	err := userCollection.FindOne(ctxDB, filter).Decode(&existingUser)

	if err != nil {
		one, errAdd := userCollection.InsertOne(ctxDB, user)
		if errAdd != nil {
			return nil, errAdd
		}
		return one, nil
	}
	return nil, nil
}
func ManageWs(s *server.Server, ctx *gin.Context) {
	res, errAdding := addUserToDb(ctx)
	result := responses.Response{
		Status:  http.StatusBadRequest,
		Message: "adding user to database not succeded",
		Data:    map[string]interface{}{"data": res},
	}
	if errAdding != nil {
		ctx.JSON(result.Status, result)
		return
	}
	ws, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := server.NewClient(ws, ctx)
	s.AddClient(client)
	client.Server = s
	go client.WritePump()
	go client.ReadPump()
}
