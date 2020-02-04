package chat

import (
	"github.com/google/uuid"
)

// Message struc to use every where
type Message struct {
	Avatar string `json:"avatar"`
	Msg    string `json:"msg"`
}

// NewMessage to publish in the chat
func NewMessage(avatar, msg string) Message {
	return Message{avatar, msg}
}

// UUID Generate a uuid
func UUID() string {
	id := uuid.New()
	return id.String()
}
