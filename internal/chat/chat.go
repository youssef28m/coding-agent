package chat

import (
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/youssef28m/codingAgent/internal/db/sqlcgen"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
	MessageRoleTool      MessageRole = "tool"
)

type Chat struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

func toDomainChat(chat sqlcgen.Chat) Chat {
	return Chat{
		ID:        chat.ID,
		Name:      chat.Name,
		CreatedAt: chat.CreatedAt,
	}
}

type Message struct {
	ID        int64
	ChatID    int64
	Role      MessageRole
	Content   string
	CreatedAt time.Time
}

func toDomainMessage(message sqlcgen.Message) Message {
	return Message{
		ID:        message.ID,
		ChatID:    message.ChatID,
		Role:      MessageRole(message.Role),
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func toOpenAIMessages(messages []Message) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, 0, len(messages))

	for _, msg := range messages {
		result = append(result, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	return result
}
