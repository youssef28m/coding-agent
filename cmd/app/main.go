package main

import (
	"log"
	"os"
)

func main() {

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY is not set")
	}

	// llmClient := llm.NewGroqLLM(apiKey, "openai/gpt-oss-20b")

}
