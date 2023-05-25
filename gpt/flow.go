package gpt

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Flow is a conversation flow. It stores a sequence of messages.
type Flow struct {
	id        uuid.UUID
	createdBy uuid.UUID
	createdAt time.Time
	// history of messages
	history []Message
	// initialSystemMessageBuilder is initial system message for chatbot, to help it pretend a role.
	// It is a stringer because it will be converted to initial system message, which is a string.
	// Is up to user how this initial system message should look like. Why not a string? Because
	// this could be some huge struct.
	initialSystemMessageBuilder fmt.Stringer

	messenger Messenger
}

// NewFlow creates a new flow.
func NewFlow(createdBy uuid.UUID, initialSystemMessage fmt.Stringer, messenger Messenger) *Flow {
	return &Flow{
		id:                          uuid.New(),
		createdAt:                   time.Now(),
		history:                     make([]Message, 0),
		createdBy:                   createdBy,
		initialSystemMessageBuilder: initialSystemMessage,
		messenger:                   messenger,
	}
}

// SendMessage sends a message and returns Response and error.
func (f *Flow) SendMessage(ctx context.Context, message string) (Response, error) {
	msg := NewMessage(f.createdBy, MessageTypeUser, message)

	payload := MessengerPayload{
		Message:                     *msg,
		InitialSystemMessageBuilder: f.initialSystemMessageBuilder,
	}

	return f.messenger.HandleMessage(ctx, payload)
}
