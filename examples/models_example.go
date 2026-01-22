// Package main demonstrates how to use the GitHub Models API integration.
// This example shows basic usage of the models package to interact with LLMs.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dreadwitdastacc-IFA/validatord/internal/models"
)

func main() {
	fmt.Println("GitHub Models API Example")
	fmt.Println("=========================")
	fmt.Println()

	// Get GitHub token from environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Error: GITHUB_TOKEN environment variable not set\n" +
			"Please create a personal access token with 'models' scope at:\n" +
			"https://github.com/settings/tokens")
	}

	// Initialize Models client with token
	m, err := models.NewWithToken(token)
	if err != nil {
		log.Fatalf("Failed to initialize models client: %v", err)
	}

	fmt.Printf("Initialized with default model: %s\n\n", m.GetDefaultModel())

	// Example 1: Simple chat with default model
	fmt.Println("Example 1: Simple chat")
	fmt.Println("----------------------")
	prompt := "What is the capital of France?"
	fmt.Printf("Prompt: %s\n", prompt)

	response, err := m.Chat(prompt)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n\n", response)
	}

	// Example 2: Chat with a specific model
	fmt.Println("Example 2: Chat with specific model")
	fmt.Println("------------------------------------")
	prompt2 := "Explain recursion in one sentence"
	modelName := "openai/gpt-4o-mini"
	fmt.Printf("Prompt: %s\n", prompt2)
	fmt.Printf("Model: %s\n", modelName)

	response2, err := m.ChatWithModel(prompt2, modelName)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n\n", response2)
	}

	// Example 3: Advanced usage with custom parameters
	fmt.Println("Example 3: Advanced request with custom parameters")
	fmt.Println("--------------------------------------------------")
	req := &models.ChatRequest{
		Model: "openai/gpt-4o-mini",
		Messages: []models.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that provides concise answers.",
			},
			{
				Role:    "user",
				Content: "What are the three primary colors?",
			},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	fmt.Println("Request:")
	fmt.Printf("  Model: %s\n", req.Model)
	fmt.Printf("  Temperature: %.1f\n", req.Temperature)
	fmt.Printf("  Max Tokens: %d\n", req.MaxTokens)
	fmt.Printf("  Messages: %d\n", len(req.Messages))

	resp, err := m.CallModel(req)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		if len(resp.Choices) > 0 {
			fmt.Printf("\nResponse: %s\n", resp.Choices[0].Message.Content)
			fmt.Printf("Model Used: %s\n", resp.Model)
		}
	}

	fmt.Println("\n=========================")
	fmt.Println("Example completed!")
}
