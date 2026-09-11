package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/youssef28m/codingAgent/internal/llm"
)

type chatService struct {
	repository Repository
	llm        llm.LLM
}

func NewChatService(repository Repository, llm llm.LLM) *chatService {
	return &chatService{
		repository: repository,
		llm:        llm,
	}
}

func (c chatService) CreateChat(ctx context.Context, name string) (Chat, error) {
	chat, err := c.repository.CreateChat(ctx, name, time.Now())
	if err != nil {
		return Chat{}, fmt.Errorf("failed to create chat: %w", err)
	}
	return chat, nil
}

func (c chatService) GetChatByID(ctx context.Context, id int64) (Chat, error) {
	chat, err := c.repository.GetChatByID(ctx, id)
	if err != nil {
		return Chat{}, fmt.Errorf("failed to get chat by ID: %w", err)
	}
	return chat, nil
}

func (c chatService) ListChats(ctx context.Context) ([]Chat, error) {
	chats, err := c.repository.ListChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list chats: %w", err)
	}
	return chats, nil
}

func (c chatService) DeleteChat(ctx context.Context, id int64) error {
	err := c.repository.DeleteChat(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}
	return nil
}

func (c chatService) SendMessage(ctx context.Context, chatID int64, content string) (Message, error) {

	// store user message
	_, err := c.repository.CreateMessage(ctx, chatID, string(MessageRoleUser), content, time.Now())
	if err != nil {
		return Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	messages, err := c.repository.ListMessagesByChatID(ctx, chatID)
	if err != nil {
		return Message{}, fmt.Errorf("failed to get message history: %w", err)
	}
	

	response, err := c.llm.GenerateResponse(ctx, toOpenAIMessages(messages))
	if err != nil {
		return Message{}, fmt.Errorf("failed to get response from llm: %w", err)
	}


	// store llm message
	result, err := c.repository.CreateMessage(
		ctx,
		chatID,
		string(MessageRoleAssistant),
		response,
		time.Now(),
	)

	if err != nil {
		return Message{}, fmt.Errorf("failed to store llm message: %w", err)
	}

	

	return result, nil
}
