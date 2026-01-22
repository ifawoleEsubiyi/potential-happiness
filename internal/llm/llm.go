// Package llm provides functionality to interact with GitHub Models API
// for running large language models. This enables easy integration with
// various AI models hosted on GitHub.
package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultAPIEndpoint is the default GitHub Models API endpoint
	DefaultAPIEndpoint = "https://models.inference.ai.azure.com"

	// DefaultModel is the default model to use if none is specified
	DefaultModel = "gpt-4o"

	// MaxTokens is the default maximum number of tokens to generate
	MaxTokens = 1000

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Message represents a chat message in the conversation
type Message struct {
	Role    string `json:"role"`    // "system", "user", or "assistant"
	Content string `json:"content"` // The message content
}

// CompletionRequest represents a request to the GitHub Models API
type CompletionRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Choice represents a single completion choice from the API
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CompletionResponse represents the response from the GitHub Models API
type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Config holds the configuration for the LLM client
type Config struct {
	APIEndpoint string
	Token       string
	Model       string
	Timeout     time.Duration
}

// Client manages interactions with the GitHub Models API.
// This type is safe for concurrent use by multiple goroutines.
type Client struct {
	mu         sync.RWMutex
	config     Config
	httpClient *http.Client
}

// New creates a new LLM client with default configuration
func New() *Client {
	return &Client{
		config: Config{
			APIEndpoint: DefaultAPIEndpoint,
			Model:       DefaultModel,
			Timeout:     DefaultTimeout,
		},
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewWithConfig creates a new LLM client with custom configuration
func NewWithConfig(cfg Config) (*Client, error) {
	if cfg.APIEndpoint == "" {
		cfg.APIEndpoint = DefaultAPIEndpoint
	}
	if cfg.Model == "" {
		cfg.Model = DefaultModel
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// SetToken sets the GitHub token for API authentication
func (c *Client) SetToken(token string) error {
	if token == "" {
		return errors.New("token cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.Token = token
	return nil
}

// SetModel sets the model to use for completions
func (c *Client) SetModel(model string) error {
	if model == "" {
		return errors.New("model cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.Model = model
	return nil
}

// GetModel returns the currently configured model
func (c *Client) GetModel() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Model
}

// GetAPIEndpoint returns the currently configured API endpoint
func (c *Client) GetAPIEndpoint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.APIEndpoint
}

// HasToken returns whether a token has been configured
func (c *Client) HasToken() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Token != ""
}

// CreateCompletion sends a completion request to the GitHub Models API
func (c *Client) CreateCompletion(messages []Message) (*CompletionResponse, error) {
	if len(messages) == 0 {
		return nil, errors.New("messages cannot be empty")
	}

	c.mu.RLock()
	token := c.config.Token
	model := c.config.Model
	endpoint := c.config.APIEndpoint
	c.mu.RUnlock()

	if token == "" {
		return nil, errors.New("GitHub token not configured")
	}

	// Construct the request
	req := CompletionRequest{
		Messages:  messages,
		Model:     model,
		MaxTokens: MaxTokens,
	}

	// Marshal request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := strings.TrimSuffix(endpoint, "/") + "/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var completionResp CompletionResponse
	if err := json.Unmarshal(respBody, &completionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &completionResp, nil
}

// SimpleCompletion is a helper method to create a completion from a single user message
func (c *Client) SimpleCompletion(prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	resp, err := c.CreateCompletion(messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatCompletion creates a completion with a system message and user message
func (c *Client) ChatCompletion(systemPrompt, userPrompt string) (string, error) {
	if userPrompt == "" {
		return "", errors.New("user prompt cannot be empty")
	}

	messages := []Message{
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	// Only add system message if provided
	if systemPrompt != "" {
		messages = append([]Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
		}, messages...)
	}

	resp, err := c.CreateCompletion(messages)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
