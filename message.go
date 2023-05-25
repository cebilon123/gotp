package gpt

import "time"

type Message struct {
	CreatedAt time.Time
	Role      int    `json:"role"`
	Content   string `json:"content"`
}

type Messages []Message

func (m Messages) toJsonMessages() []conversationMessageJson {
	jsonMessages := make([]conversationMessageJson, len(m))
	for i, message := range m {
		jsonMessages[i] = conversationMessageJson{
			Role:    MessageTypeToString(message.Role),
			Content: message.Content,
		}
	}
	return jsonMessages
}
