package models

type RoomDB struct {
	id       string        `bson:"id"`
	name     string        `bson:"name"`
	users    []string      `bson:"users"`
	messages []interface{} `bson:"messages"`
}
