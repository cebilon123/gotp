package gpt

import (
	"github.com/google/uuid"
	"time"
)

// MessageType is the type of the message. It can be system, user or assistant.
type MessageType int

func (m MessageType) String() string {
	switch m {
	case MessageTypeSystem:
		return "system"
	case MessageTypeUser:
		return "user"
	case MessageTypeAssistant:
		return "assistant"
	default:
		return "user"
	}
}

const (
	MessageTypeSystem MessageType = iota
	MessageTypeUser
	MessageTypeAssistant
)

// Message is a message in a conversation.
type Message struct {
	id          uuid.UUID
	createdAt   time.Time
	createBy    uuid.UUID
	messageType MessageType
	content     string
}

func NewMessage(createBy uuid.UUID, messageType MessageType, content string) *Message {
	return &Message{
		id:          uuid.New(),
		createdAt:   time.Now(),
		createBy:    createBy,
		messageType: messageType,
		content:     content,
	}
}
