package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/youssef28m/codingAgent/internal/chat"
	"github.com/youssef28m/codingAgent/internal/db"
	"github.com/youssef28m/codingAgent/internal/llm"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY is not set")
	}

	ctx := context.Background()

	conn, err := db.NewSQLiteDB("app.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	repository := chat.NewRepository(conn)
	llmClient := llm.NewGroqLLM(apiKey, "openai/gpt-oss-20b")
	chatService := chat.NewChatService(repository, llmClient)

	chat1, err := chatService.CreateChat(ctx, "My first chat")
	if err != nil {
		log.Fatalf("failed to create a chat :%v", err)
	}
	message, err := chatService.SendMessage(ctx, chat1.ID, "Hello, how are you?")

	fmt.Printf("%q: %q", message.Role, message.Content)

}
