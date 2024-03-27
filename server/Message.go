package server

import (
	"encoding/json"
	"log"
)

type Message struct {
	Message     string `json:"message" bson:"message"`
	Sender      string `json:"sender" bson:"sender"`
	Destination string `json:"destination" bson:"destination"`
	Options     string `json:"options" bson:"options"`
}

func ByteToMessage(read []byte) *Message {
	mess := Message{}
	json.Unmarshal(read, &mess)
	return &mess
}

func MessageToByte(read *Message) *[]byte {
	byteMessage, err := json.Marshal(read)
	if err != nil {
		log.Fatal(err)
	}
	return &byteMessage
}
