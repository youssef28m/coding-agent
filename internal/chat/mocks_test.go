package chat

import (
	"context"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
)

// CreateChat(ctx context.Context, name string, createdAt time.Time) (Chat, error)
// GetChatByID(ctx context.Context, id int64) (Chat, error)
// ListChats(ctx context.Context) ([]Chat, error)
// DeleteChat(ctx context.Context, id int64) error

// CreateMessage(ctx context.Context, chatID int64, role string, content string, createdAt time.Time) (Message, error)
// ListMessagesByChatID(ctx context.Context, chatID int64) ([]Message, error)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateChat(ctx context.Context, name string, createdAt time.Time) (Chat, error) {
	args := m.Called(ctx, name, createdAt)
	return args.Get(0).(Chat), args.Error(1)
}

func (m *MockRepository) GetChatByID(ctx context.Context, id int64) (Chat, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Chat), args.Error(1)
}

func (m *MockRepository) ListChats(ctx context.Context) ([]Chat, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Chat), args.Error(1)
}

func (m *MockRepository) DeleteChat(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateMessage(ctx context.Context, chatID int64, role, content string, createdAt time.Time) (Message, error) {
	args := m.Called(ctx, chatID, role, content, createdAt)
	return args.Get(0).(Message), args.Error(1)
}

func (m *MockRepository) ListMessagesByChatID(ctx context.Context, chatID int64) ([]Message, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).([]Message), args.Error(1)
}

type MockLLM struct {
	mock.Mock
}

func (m *MockLLM) GenerateResponse(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	args := m.Called(ctx, messages)
	return args.String(0), args.Error(1)
}
