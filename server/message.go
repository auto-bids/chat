package server

import (
	"encoding/json"
	"log"
)

type Message struct {
	Message     string `json:"message"`
	Sender      string `json:"sender"`
	Destination string `json:"destination"`
	Options     string `json:"options"`
}

func ByteToMessage(read []byte) Message {
	mess := Message{}
	json.Unmarshal(read, &mess)
	return mess
}
func messageToByte(read *Message) []byte {
	byteMessage, err := json.Marshal(read)
	if err != nil {
		log.Fatal(err)
	}
	return byteMessage
}
