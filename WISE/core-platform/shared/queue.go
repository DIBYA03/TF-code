package shared

import (
	"context"
)

type Message struct {
	ID      string  `json:"id"`
	GroupID *string `json:"groupId"`
	Body    []byte  `json:"body"`
	Delay   *int64  `json:"delay"`
}

type MessageHandler interface {
	HandleMessage(context.Context, Message) error
}

type MessageQueue interface {
	ReceiveMessages(context.Context, MessageHandler) error
	SendMessages([]Message) (SendMessageResult, error)
	URL() *string
}

type SendMessageResult struct {
	SuccessIDs []string
	FailedIDs  []string
}
