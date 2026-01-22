// Package models provides GitHub Models API integration for the validator daemon,
// allowing easy access to LLM capabilities via GitHub's AI inference API.
package models

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
	DefaultAPIEndpoint = "https://models.github.ai/inference/chat/completions"

	// MaxPromptLength is the maximum allowed prompt length
	MaxPromptLength = 100000

	// MaxResponseTokens is the default maximum response tokens
	MaxResponseTokens = 4096

	// DefaultModel is the default model to use if none specified
	DefaultModel = "openai/gpt-4o-mini"

	// DefaultTimeout is the default HTTP request timeout
	DefaultTimeout = 60 * time.Second
)

// Message represents a chat message in the conversation
type Message struct {
	Role    string `json:"role"`    // Role of the message sender (user, assistant, system)
	Content string `json:"content"` // Content of the message
}

// ChatRequest represents a request to the GitHub Models chat completion API
type ChatRequest struct {
	Model       string    `json:"model"`                 // Model to use for inference
	Messages    []Message `json:"messages"`              // Messages in the conversation
	MaxTokens   int       `json:"max_tokens,omitempty"`  // Maximum tokens to generate
	Temperature float64   `json:"temperature,omitempty"` // Sampling temperature (0.0 to 2.0)
}

// ChatChoice represents a single completion choice from the API
type ChatChoice struct {
	Message Message `json:"message"`
	Index   int     `json:"index"`
}

// ChatResponse represents the response from the GitHub Models API
type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
	Model   string       `json:"model"`
	ID      string       `json:"id"`
}

// ModelConfig holds the configuration for GitHub Models API operations
type ModelConfig struct {
	Token        string        // GitHub personal access token with 'models' scope
	Endpoint     string        // API endpoint URL
	DefaultModel string        // Default model to use
	Timeout      time.Duration // HTTP request timeout
}

// Models manages GitHub Models API operations.
// This type is safe for concurrent use by multiple goroutines.
type Models struct {
	mu     sync.RWMutex
	config ModelConfig
	client *http.Client
}

// New creates a new Models instance with default configuration.
// Note: A GitHub personal access token with 'models' scope must be set
// using SetToken() before making API calls.
func New() *Models {
	return &Models{
		config: ModelConfig{
			Token:        "",
			Endpoint:     DefaultAPIEndpoint,
			DefaultModel: DefaultModel,
			Timeout:      DefaultTimeout,
		},
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewWithToken creates a new Models instance with the specified GitHub token.
// The token must have the 'models' scope.
func NewWithToken(token string) (*Models, error) {
	if err := ValidateToken(token); err != nil {
		return nil, err
	}

	m := New()
	m.config.Token = token
	return m, nil
}

// ValidateToken validates that a GitHub token is properly formatted.
func ValidateToken(token string) error {
	if token == "" {
		return errors.New("GitHub token cannot be empty")
	}

	// GitHub tokens typically start with specific prefixes
	token = strings.TrimSpace(token)
	if len(token) < 10 {
		return errors.New("GitHub token appears to be invalid (too short)")
	}

	return nil
}

// SetToken sets the GitHub personal access token for API authentication
func (m *Models) SetToken(token string) error {
	if err := ValidateToken(token); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Token = token
	return nil
}

// GetEndpoint returns the configured API endpoint
func (m *Models) GetEndpoint() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Endpoint
}

// GetDefaultModel returns the configured default model
func (m *Models) GetDefaultModel() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.DefaultModel
}

// SetDefaultModel sets the default model to use for API calls
func (m *Models) SetDefaultModel(model string) error {
	if model == "" {
		return errors.New("model name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.DefaultModel = model
	return nil
}

// HasToken returns whether a token has been configured
func (m *Models) HasToken() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Token != ""
}

// ValidateChatRequest validates a chat request
func ValidateChatRequest(req *ChatRequest) error {
	if req == nil {
		return errors.New("chat request cannot be nil")
	}

	if req.Model == "" {
		return errors.New("model must be specified")
	}

	if len(req.Messages) == 0 {
		return errors.New("messages cannot be empty")
	}

	// Validate each message
	for i, msg := range req.Messages {
		if msg.Role == "" {
			return fmt.Errorf("message %d: role cannot be empty", i)
		}
		if msg.Content == "" {
			return fmt.Errorf("message %d: content cannot be empty", i)
		}
		if len(msg.Content) > MaxPromptLength {
			return fmt.Errorf("message %d: content exceeds maximum length of %d", i, MaxPromptLength)
		}
		// Validate role is one of the expected values
		if msg.Role != "user" && msg.Role != "assistant" && msg.Role != "system" {
			return fmt.Errorf("message %d: invalid role '%s' (must be 'user', 'assistant', or 'system')", i, msg.Role)
		}
	}

	// Validate temperature if specified
	if req.Temperature < 0.0 || req.Temperature > 2.0 {
		return errors.New("temperature must be between 0.0 and 2.0")
	}

	// Validate max tokens if specified
	if req.MaxTokens < 0 {
		return errors.New("max_tokens must be non-negative")
	}

	return nil
}

// CallModel makes a call to the GitHub Models API with the specified request.
// Returns the API response or an error if the request fails.
func (m *Models) CallModel(req *ChatRequest) (*ChatResponse, error) {
	if err := ValidateChatRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	m.mu.RLock()
	token := m.config.Token
	endpoint := m.config.Endpoint
	m.mu.RUnlock()

	if token == "" {
		return nil, errors.New("GitHub token not configured")
	}

	// Marshal request to JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	httpReq.Header.Set("Accept", "application/vnd.github+json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := m.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &chatResp, nil
}

// Chat is a convenience method to send a single user message and get a response.
// It uses the default model configured for this instance.
func (m *Models) Chat(prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	m.mu.RLock()
	defaultModel := m.config.DefaultModel
	m.mu.RUnlock()

	req := &ChatRequest{
		Model: defaultModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	resp, err := m.CallModel(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices returned from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatWithModel is a convenience method to send a single user message with a specific model.
func (m *Models) ChatWithModel(prompt, model string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}
	if model == "" {
		return "", errors.New("model cannot be empty")
	}

	req := &ChatRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	resp, err := m.CallModel(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices returned from API")
	}

	return resp.Choices[0].Message.Content, nil
}
