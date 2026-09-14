package chat

import (
	"errors"
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

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrCreateChat        = errors.New("failed to create chat")
	ErrGetChatByID       = errors.New("failed to get chat by ID")
	ErrListChats         = errors.New("failed to list chats")
	ErrDeleteChat        = errors.New("failed to delete chat")
	ErrCreateMessage     = errors.New("failed to create message")
	ErrGetMessageHistory = errors.New("failed to get message history")
	ErrGetLLMResponse    = errors.New("failed to get response from llm")
	ErrStoreLLMMessage   = errors.New("failed to store llm message")
)
