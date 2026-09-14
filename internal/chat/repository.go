package chat

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/youssef28m/codingAgent/internal/db/sqlcgen"
)

type Repository interface {
	CreateChat(ctx context.Context, name string, createdAt time.Time) (Chat, error)
	GetChatByID(ctx context.Context, id int64) (Chat, error)
	ListChats(ctx context.Context) ([]Chat, error)
	DeleteChat(ctx context.Context, id int64) error

	CreateMessage(ctx context.Context, chatID int64, role string, content string, createdAt time.Time) (Message, error)
	ListMessagesByChatID(ctx context.Context, chatID int64) ([]Message, error)
}

type repository struct {
	q *sqlcgen.Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		q: sqlcgen.New(db),
	}
}

func (r *repository) CreateChat(ctx context.Context, name string, createdAt time.Time) (Chat, error) {
	c, err := r.q.CreateChat(ctx, sqlcgen.CreateChatParams{
		Name:      name,
		CreatedAt: createdAt,
	})
	if err != nil {
		return Chat{}, err
	}
	return toDomainChat(c), nil
}

func (r *repository) GetChatByID(ctx context.Context, id int64) (Chat, error) {
	c, err := r.q.GetChatByID(ctx, id)
	if err != nil {
		return Chat{}, fmt.Errorf("failed to get chat by ID: %w", err)
	}

	return toDomainChat(c), nil
}

func (r *repository) ListChats(ctx context.Context) ([]Chat, error) {
	c, err := r.q.ListChats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list chats: %w", err)
	}

	result := make([]Chat, 0, len(c))
	for _, value := range c {
		result = append(result, toDomainChat(value))
	}

	return result, nil
}

func (r *repository) DeleteChat(ctx context.Context, id int64) error {
	err := r.q.DeleteChat(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (r *repository) CreateMessage(ctx context.Context, chatID int64, role string, content string, createdAt time.Time) (Message, error) {
	m, err := r.q.CreateMessage(ctx, sqlcgen.CreateMessageParams{
		ChatID:    chatID,
		Role:      role,
		Content:   content,
		CreatedAt: createdAt,
	})
	if err != nil {
		return Message{}, err
	}
	return toDomainMessage(m), nil
}

func (r *repository) ListMessagesByChatID(ctx context.Context, chatID int64) ([]Message, error) {
	messages, err := r.q.ListMessagesByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	domainMessages := make([]Message, len(messages))
	for i, message := range messages {
		domainMessages[i] = toDomainMessage(message)
	}
	return domainMessages, nil
}
