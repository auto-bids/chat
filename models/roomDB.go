package models

type RoomDB struct {
	Id       string      `bson:"Id"`
	Name     string      `bson:"Name"`
	Users    []string    `bson:"Users"`
	Messages []MessageDB `bson:"Messages"`
}
