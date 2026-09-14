package chat

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateChat_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLLM := new(MockLLM)
	service := NewChatService(mockRepo, mockLLM)

	expectedChat := Chat{ID: 1, Name: "my chat"}

	mockRepo.On("CreateChat", mock.Anything, "my chat", mock.Anything).
		Return(expectedChat, nil)

	result, err := service.CreateChat(context.Background(), "my chat")

	assert.NoError(t, err)
	assert.Equal(t, expectedChat, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateChat_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLLM := new(MockLLM)
	service := NewChatService(mockRepo, mockLLM)
	repositoryErr := errors.New("db connection failed")

	mockRepo.On("CreateChat", mock.Anything, "my chat", mock.Anything).
		Return(Chat{}, repositoryErr)

	result, err := service.CreateChat(context.Background(), "my chat")

	assert.ErrorIs(t, err, ErrCreateChat)
	assert.ErrorIs(t, err, repositoryErr)
	assert.Equal(t, Chat{}, result)
	mockRepo.AssertExpectations(t)
}

func TestSendMessage_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLLM := new(MockLLM)
	service := NewChatService(mockRepo, mockLLM)

	chatID := int64(1)
	userContent := "hello"
	aiReply := "hi, how can I help?"

	userMsg := Message{ID: 10, ChatID: chatID, Role: "user", Content: userContent}
	history := []Message{userMsg}
	finalMsg := Message{ID: 11, ChatID: chatID, Role: "assistant", Content: aiReply}

	// 1st call: store the user's message
	mockRepo.On("CreateMessage", mock.Anything, chatID, "user", userContent, mock.Anything).
		Return(userMsg, nil)

	// 2nd call: fetch history to send to the LLM
	mockRepo.On("ListMessagesByChatID", mock.Anything, chatID).
		Return(history, nil)

	// 3rd call: LLM generates a reply
	mockLLM.On("GenerateResponse", mock.Anything, mock.Anything).
		Return(aiReply, nil)

	// 4th call: store the assistant's reply
	mockRepo.On("CreateMessage", mock.Anything, chatID, "assistant", aiReply, mock.Anything).
		Return(finalMsg, nil)

	result, err := service.SendMessage(context.Background(), chatID, userContent)

	assert.NoError(t, err)
	assert.Equal(t, finalMsg, result)
	mockRepo.AssertExpectations(t)
	mockLLM.AssertExpectations(t)
}

func TestSendMessage_Failures(t *testing.T) {
	chatID := int64(1)
	userContent := "hello"
	userMsg := Message{ID: 10, ChatID: chatID, Role: "user", Content: userContent}

	tests := []struct {
		name       string
		setupMocks func(repo *MockRepository, llm *MockLLM)
		wantErr    error
	}{
		{
			name: "store user message fails",
			setupMocks: func(repo *MockRepository, llm *MockLLM) {
				repo.On("CreateMessage", mock.Anything, chatID, "user", userContent, mock.Anything).
					Return(Message{}, errors.New("db down"))
			},
			wantErr: ErrCreateMessage,
		},
		{
			name: "message history lookup fails",
			setupMocks: func(repo *MockRepository, llm *MockLLM) {
				repo.On("CreateMessage", mock.Anything, chatID, "user", userContent, mock.Anything).
					Return(userMsg, nil)
				repo.On("ListMessagesByChatID", mock.Anything, chatID).
					Return([]Message{}, errors.New("history unavailable"))
			},
			wantErr: ErrGetMessageHistory,
		},
		{
			name: "llm call fails",
			setupMocks: func(repo *MockRepository, llm *MockLLM) {
				repo.On("CreateMessage", mock.Anything, chatID, "user", userContent, mock.Anything).
					Return(userMsg, nil)
				repo.On("ListMessagesByChatID", mock.Anything, chatID).
					Return([]Message{userMsg}, nil)
				llm.On("GenerateResponse", mock.Anything, mock.Anything).
					Return("", errors.New("LLM timeout"))
			},
			wantErr: ErrGetLLMResponse,
		},
		{
			name: "assistant message store fails",
			setupMocks: func(repo *MockRepository, llm *MockLLM) {
				aiReply := "hi, how can I help?"
				repo.On("CreateMessage", mock.Anything, chatID, "user", userContent, mock.Anything).
					Return(userMsg, nil)
				repo.On("ListMessagesByChatID", mock.Anything, chatID).
					Return([]Message{userMsg}, nil)
				llm.On("GenerateResponse", mock.Anything, mock.Anything).
					Return(aiReply, nil)
				repo.On("CreateMessage", mock.Anything, chatID, "assistant", aiReply, mock.Anything).
					Return(Message{}, errors.New("db write failed"))
			},
			wantErr: ErrStoreLLMMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockLLM := new(MockLLM)
			tt.setupMocks(mockRepo, mockLLM)

			service := NewChatService(mockRepo, mockLLM)
			result, err := service.SendMessage(context.Background(), chatID, userContent)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, Message{}, result)
			mockRepo.AssertExpectations(t)
			mockLLM.AssertExpectations(t)
		})
	}
}

func TestGetChatByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLLm := new(MockLLM)
	service := NewChatService(mockRepo, mockLLm)

	expectedChat := Chat{ID: 2, Name: "chat"}

	id := int64(2)

	mockRepo.On("GetChatByID", mock.Anything, id).
		Return(expectedChat, nil)

	result, err := service.GetChatByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expectedChat, result)
	mockRepo.AssertExpectations(t)
}

func TestGetChatByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLLm := new(MockLLM)
	service := NewChatService(mockRepo, mockLLm)

	id := int64(2)

	mockRepo.On("GetChatByID", mock.Anything, id).
		Return(Chat{}, sql.ErrNoRows)

	result, err := service.GetChatByID(context.Background(), id)

	assert.ErrorIs(t, err, ErrChatNotFound)
	assert.Equal(t, Chat{}, result)
	mockRepo.AssertExpectations(t)

}

func TestGetChatByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewChatService(mockRepo, new(MockLLM))
	repositoryErr := errors.New("database unavailable")
	id := int64(2)

	mockRepo.On("GetChatByID", mock.Anything, id).
		Return(Chat{}, repositoryErr)

	result, err := service.GetChatByID(context.Background(), id)

	assert.ErrorIs(t, err, ErrGetChatByID)
	assert.ErrorIs(t, err, repositoryErr)
	assert.Equal(t, Chat{}, result)
	mockRepo.AssertExpectations(t)
}

func TestListChats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewChatService(mockRepo, new(MockLLM))
		wantChats := []Chat{{ID: 1, Name: "first"}, {ID: 2, Name: "second"}}

		mockRepo.On("ListChats", mock.Anything).Return(wantChats, nil)

		gotChats, err := service.ListChats(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, wantChats, gotChats)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewChatService(mockRepo, new(MockLLM))
		repositoryErr := errors.New("database unavailable")

		mockRepo.On("ListChats", mock.Anything).Return([]Chat(nil), repositoryErr)

		gotChats, err := service.ListChats(context.Background())

		assert.ErrorIs(t, err, ErrListChats)
		assert.ErrorIs(t, err, repositoryErr)
		assert.Nil(t, gotChats)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteChat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewChatService(mockRepo, new(MockLLM))
		id := int64(1)

		mockRepo.On("DeleteChat", mock.Anything, id).Return(nil)

		err := service.DeleteChat(context.Background(), id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewChatService(mockRepo, new(MockLLM))
		repositoryErr := errors.New("database unavailable")
		id := int64(1)

		mockRepo.On("DeleteChat", mock.Anything, id).Return(repositoryErr)

		err := service.DeleteChat(context.Background(), id)

		assert.ErrorIs(t, err, ErrDeleteChat)
		assert.ErrorIs(t, err, repositoryErr)
		mockRepo.AssertExpectations(t)
	})
}
