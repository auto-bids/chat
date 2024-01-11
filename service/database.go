package service

import (
	"chat/server"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(os.Getenv("DB"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB")))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func sendMessage(conn *mongo.Client, roomId string, message *server.Message) *mongo.InsertOneResult {
	one, err := conn.Database("chat").Collection("messages").InsertOne(context.TODO(), message)
	if err != nil {
		log.Fatal(err)
	}
	return one
}
func addRoom(conn *mongo.Client, roomId string) *mongo.InsertOneResult {
	one, err := conn.Database("chat").Collection("room").InsertOne(context.TODO(), roomId)
	if err != nil {
		log.Fatal(err)
	}
	return one
}
func getMessages(conn *mongo.Client, roomId string, message *server.Message) *mongo.InsertOneResult {

}
