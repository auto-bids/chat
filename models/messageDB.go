package models

type MessageDB struct {
	Sender  string `bson:"Sender" json:"Sender"`
	Message string `bson:"Message" json:"Message"`
	Time    int64  `bson:"Time"`
}
type MessageUnwindDB struct {
	Messages MessageDB `bson:"messages" json:"messages"`
}
