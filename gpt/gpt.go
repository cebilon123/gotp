package gpt

import (
	"context"
	"fmt"
)

type Response struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type MessengerPayload struct {
	Message                     Message
	InitialSystemMessageBuilder fmt.Stringer
}

// Messenger should be implemented by any type that wants to handle messages.
// It should be goroutine safe, because as the intention there should be only
// one instance of Messenger used.
type Messenger interface {
	HandleMessage(ctx context.Context, payload MessengerPayload) (Response, error)
}

type BasicMessenger struct {
}

func (b *BasicMessenger) HandleMessage(ctx context.Context, payload MessengerPayload) (Response, error) {
	//TODO implement me
	panic("implement me")
}
