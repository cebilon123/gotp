package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Conversation struct {
	createdAt time.Time
	history   Messages
	client    *http.Client
	chatType  string
	apiKey    string
}

func NewConversation(client *http.Client, chatType string, apiKey string) Conversation {
	return Conversation{
		createdAt: time.Now(),
		history:   make([]Message, 0),
		client:    client,
		chatType:  chatType,
		apiKey:    apiKey,
	}
}

func (c *Conversation) AddSystemMessage(message string) {
	c.history = append(c.history, Message{
		CreatedAt: time.Now(),
		Role:      MessageTypeSystem,
		Content:   message,
	})
}

func (c *Conversation) AddUserMessage(message string) {
	c.history = append(c.history, Message{
		CreatedAt: time.Now(),
		Role:      MessageTypeUser,
		Content:   message,
	})
}

func (c *Conversation) SendMessage(ctx context.Context, message Message) (*Message, error) {
	res, err := c.sendMessageUsingHttp(ctx, message)
	if err != nil {
		return nil, err
	}

	msg := Message{
		CreatedAt: time.Now(),
		Role:      StringToMessageType(res.Choices[0].Message.Role),
		Content:   res.Choices[0].Message.Content,
	}
	c.history = append(c.history, message)
	c.history = append(c.history, msg)

	return &msg, nil
}

type conversationJson struct {
	Model    string                    `json:"model"`
	Messages []conversationMessageJson `json:"messages"`
}

type conversationMessageJson struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const endpoint = "https://api.openai.com/v1/chat/completions"

type ResponseJson struct {
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

func (c *Conversation) sendMessageUsingHttp(ctx context.Context, message Message) (*ResponseJson, error) {
	conv := &conversationJson{
		Model:    c.chatType,
		Messages: c.history.toJsonMessages(),
	}
	conv.Messages = append(conv.Messages, conversationMessageJson{
		Role:    MessageTypeToString(message.Role),
		Content: message.Content,
	})
	body, err := json.Marshal(conv)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		resBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error status code: %d; %w", res.StatusCode, errors.Join(err, errors.New(string(resBytes))))
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responseJson ResponseJson
	err = json.Unmarshal(resBytes, &responseJson)
	if err != nil {
		return nil, err
	}

	return &responseJson, nil
}

const (
	MessageTypeSystem = iota
	MessageTypeUser
	MessageTypeAssistant
)

func MessageTypeToString(messageType int) string {
	switch messageType {
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

func StringToMessageType(messageType string) int {
	switch messageType {
	case "system":
		return MessageTypeSystem
	case "user":
		return MessageTypeUser
	case "assistant":
		return MessageTypeAssistant
	default:
		return MessageTypeUser
	}
}

type Message struct {
	CreatedAt time.Time
	Role      int    `json:"role"`
	Content   string `json:"content"`
}
