package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func connectDB() {
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
}
