// Package main demonstrates how to use the GitHub Models API integration in validatord
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dreadwitdastacc-IFA/validatord/internal/llm"
)

func main() {
	// Example 1: Create a basic LLM client
	client := llm.New()
	fmt.Println("Example 1: Basic LLM Client")
	fmt.Printf("API Endpoint: %s\n", client.GetAPIEndpoint())
	fmt.Printf("Model: %s\n", client.GetModel())
	fmt.Println()

	// Example 2: Configure with GitHub token (from environment variable)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Note: Set GITHUB_TOKEN environment variable to run API calls")
		fmt.Println("Example: export GITHUB_TOKEN=your_github_token")
		fmt.Println()
	} else {
		err := client.SetToken(token)
		if err != nil {
			log.Fatalf("Failed to set token: %v", err)
		}
		fmt.Printf("Token configured: %v\n", client.HasToken())

		// Example 3: Simple completion
		fmt.Println("\nExample 2: Simple Completion")
		response, err := client.SimpleCompletion("What are the benefits of blockchain validation?")
		if err != nil {
			log.Fatalf("Simple completion failed: %v", err)
		}
		fmt.Printf("Q: What are the benefits of blockchain validation?\n")
		fmt.Printf("A: %s\n\n", response)

		// Example 4: Chat completion with system context
		fmt.Println("Example 3: Chat Completion with System Context")
		response, err = client.ChatCompletion(
			"You are a blockchain expert who explains concepts clearly and concisely.",
			"Explain BLS signatures in 2-3 sentences",
		)
		if err != nil {
			log.Fatalf("Chat completion failed: %v", err)
		}
		fmt.Printf("Q: Explain BLS signatures\n")
		fmt.Printf("A: %s\n\n", response)

		// Example 5: Advanced usage with custom messages
		fmt.Println("Example 4: Advanced Usage with Custom Messages")
		messages := []llm.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant specialized in cryptocurrency and blockchain technology.",
			},
			{
				Role:    "user",
				Content: "What is the difference between Proof of Work and Proof of Stake?",
			},
		}

		fullResponse, err := client.CreateCompletion(messages)
		if err != nil {
			log.Fatalf("Create completion failed: %v", err)
		}

		fmt.Printf("Q: What is the difference between Proof of Work and Proof of Stake?\n")
		fmt.Printf("A: %s\n", fullResponse.Choices[0].Message.Content)
		fmt.Printf("\nToken usage:\n")
		fmt.Printf("  Prompt tokens: %d\n", fullResponse.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", fullResponse.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", fullResponse.Usage.TotalTokens)
	}

	// Example 6: Custom configuration
	fmt.Println("\nExample 5: Custom Configuration")
	customClient, err := llm.NewWithConfig(llm.Config{
		APIEndpoint: "https://models.inference.ai.azure.com",
		Model:       "gpt-4o-mini",
		Token:       token,
	})
	if err != nil {
		log.Fatalf("Failed to create custom client: %v", err)
	}
	fmt.Printf("Custom client created with model: %s\n", customClient.GetModel())

	fmt.Println("\nExamples completed successfully!")
}
