package models

type RoomDB struct {
	Id       string      `bson:"_id" json:"_id"`
	Name     string      `bson:"name" json:"name"`
	Users    []string    `bson:"users" json:"users"`
	Messages []MessageDB `bson:"messages" json:"messages"`
}

type PostPrivateRoomDB struct {
	Name  string   `bson:"name" json:"name"`
	Users []string `bson:"users" json:"users" validate:"required,min=1,max=1,email"`
}

type PostRoomDB struct {
	Name     string      `bson:"name" json:"name"`
	Users    []string    `bson:"users" json:"users" validate:"required,min=1,max=1,email"`
	Messages []MessageDB `bson:"messages" json:"messages"`
}
