package models

type UserDB struct {
	Id       string   `json:"_id" bson:"_id"`
	Username string   `json:"username" bson:"username"`
	Email    string   `json:"email" bson:"email"`
	Rooms    []string `json:"rooms" bson:"rooms"`
}
type PostUserDB struct {
	Username string   `json:"username" bson:"username" validate:"required"`
	Email    string   `json:"email" bson:"email" validate:"required,email"`
	Rooms    []string `json:"rooms" bson:"rooms" validate:"required"`
}
