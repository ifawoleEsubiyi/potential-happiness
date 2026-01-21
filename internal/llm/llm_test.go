package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("New() returned nil")
	}

	if client.GetAPIEndpoint() != DefaultAPIEndpoint {
		t.Errorf("New() APIEndpoint = %q, want %q", client.GetAPIEndpoint(), DefaultAPIEndpoint)
	}

	if client.GetModel() != DefaultModel {
		t.Errorf("New() Model = %q, want %q", client.GetModel(), DefaultModel)
	}

	if client.HasToken() {
		t.Error("New() HasToken = true, want false")
	}
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				APIEndpoint: "https://custom.endpoint.com",
				Token:       "test-token",
				Model:       "gpt-3.5-turbo",
				Timeout:     10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "config with defaults",
			config: Config{
				Token: "test-token",
			},
			wantErr: false,
		},
		{
			name:    "empty config uses defaults",
			config:  Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewWithConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewWithConfig() returned nil for valid config")
			}
		})
	}
}

func TestSetToken(t *testing.T) {
	client := New()

	// Test setting a valid token
	err := client.SetToken("test-token-123")
	if err != nil {
		t.Errorf("SetToken() error = %v, want nil", err)
	}

	if !client.HasToken() {
		t.Error("SetToken() HasToken = false, want true")
	}

	// Test setting an empty token
	err = client.SetToken("")
	if err == nil {
		t.Error("SetToken() with empty token should return error")
	}
}

func TestSetModel(t *testing.T) {
	client := New()

	// Test setting a valid model
	err := client.SetModel("gpt-3.5-turbo")
	if err != nil {
		t.Errorf("SetModel() error = %v, want nil", err)
	}

	if client.GetModel() != "gpt-3.5-turbo" {
		t.Errorf("SetModel() Model = %q, want %q", client.GetModel(), "gpt-3.5-turbo")
	}

	// Test setting an empty model
	err = client.SetModel("")
	if err == nil {
		t.Error("SetModel() with empty model should return error")
	}
}

func TestCreateCompletion(t *testing.T) {
	// Create a mock server
	mockResponse := CompletionResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-4o",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "Hello! How can I help you today?",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization: Bearer test-token, got %s", authHeader)
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server
	client, err := NewWithConfig(Config{
		APIEndpoint: server.URL,
		Token:       "test-token",
		Model:       "gpt-4o",
	})
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}

	// Test successful completion
	messages := []Message{
		{
			Role:    "user",
			Content: "Hello, how are you?",
		},
	}

	resp, err := client.CreateCompletion(messages)
	if err != nil {
		t.Errorf("CreateCompletion() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("CreateCompletion() returned nil response")
	}

	if len(resp.Choices) == 0 {
		t.Error("CreateCompletion() returned no choices")
	}

	if resp.Choices[0].Message.Content != "Hello! How can I help you today?" {
		t.Errorf("CreateCompletion() content = %q, want %q",
			resp.Choices[0].Message.Content, "Hello! How can I help you today?")
	}
}

func TestCreateCompletionErrors(t *testing.T) {
	client := New()

	// Test with no token
	messages := []Message{
		{Role: "user", Content: "test"},
	}
	_, err := client.CreateCompletion(messages)
	if err == nil {
		t.Error("CreateCompletion() without token should return error")
	}

	// Test with empty messages
	client.SetToken("test-token")
	_, err = client.CreateCompletion([]Message{})
	if err == nil {
		t.Error("CreateCompletion() with empty messages should return error")
	}
}

func TestSimpleCompletion(t *testing.T) {
	mockResponse := CompletionResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-4o",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "This is a test response",
				},
				FinishReason: "stop",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client, err := NewWithConfig(Config{
		APIEndpoint: server.URL,
		Token:       "test-token",
		Model:       "gpt-4o",
	})
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}

	// Test successful simple completion
	content, err := client.SimpleCompletion("Test prompt")
	if err != nil {
		t.Errorf("SimpleCompletion() error = %v, want nil", err)
	}

	if content != "This is a test response" {
		t.Errorf("SimpleCompletion() content = %q, want %q", content, "This is a test response")
	}

	// Test with empty prompt
	_, err = client.SimpleCompletion("")
	if err == nil {
		t.Error("SimpleCompletion() with empty prompt should return error")
	}
}

func TestChatCompletion(t *testing.T) {
	mockResponse := CompletionResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-4o",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "Chat response",
				},
				FinishReason: "stop",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client, err := NewWithConfig(Config{
		APIEndpoint: server.URL,
		Token:       "test-token",
		Model:       "gpt-4o",
	})
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}

	// Test with system and user prompts
	content, err := client.ChatCompletion("You are a helpful assistant", "Hello")
	if err != nil {
		t.Errorf("ChatCompletion() error = %v, want nil", err)
	}

	if content != "Chat response" {
		t.Errorf("ChatCompletion() content = %q, want %q", content, "Chat response")
	}

	// Test with only user prompt
	content, err = client.ChatCompletion("", "Hello")
	if err != nil {
		t.Errorf("ChatCompletion() with empty system prompt error = %v, want nil", err)
	}

	// Test with empty user prompt
	_, err = client.ChatCompletion("System prompt", "")
	if err == nil {
		t.Error("ChatCompletion() with empty user prompt should return error")
	}
}

func TestCreateCompletionHTTPError(t *testing.T) {
	// Create a server that returns an error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid token"}`))
	}))
	defer server.Close()

	client, err := NewWithConfig(Config{
		APIEndpoint: server.URL,
		Token:       "invalid-token",
		Model:       "gpt-4o",
	})
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}

	messages := []Message{
		{Role: "user", Content: "test"},
	}

	_, err = client.CreateCompletion(messages)
	if err == nil {
		t.Error("CreateCompletion() with invalid token should return error")
	}
}

func TestCreateCompletionNoChoices(t *testing.T) {
	// Create a server that returns a response with no choices
	mockResponse := CompletionResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-4o",
		Choices: []Choice{}, // Empty choices
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client, err := NewWithConfig(Config{
		APIEndpoint: server.URL,
		Token:       "test-token",
		Model:       "gpt-4o",
	})
	if err != nil {
		t.Fatalf("NewWithConfig() error = %v", err)
	}

	_, err = client.SimpleCompletion("test")
	if err == nil {
		t.Error("SimpleCompletion() with no choices should return error")
	}
}

func TestMessageStruct(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "test message",
	}

	// Test JSON marshaling
	data, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("Failed to marshal Message: %v", err)
	}

	// Test JSON unmarshaling
	var decoded Message
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Errorf("Failed to unmarshal Message: %v", err)
	}

	if decoded.Role != msg.Role {
		t.Errorf("Decoded Role = %q, want %q", decoded.Role, msg.Role)
	}

	if decoded.Content != msg.Content {
		t.Errorf("Decoded Content = %q, want %q", decoded.Content, msg.Content)
	}
}

func TestCompletionRequestStruct(t *testing.T) {
	req := CompletionRequest{
		Messages: []Message{
			{Role: "user", Content: "test"},
		},
		Model:       "gpt-4o",
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	// Test JSON marshaling
	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal CompletionRequest: %v", err)
	}

	// Test JSON unmarshaling
	var decoded CompletionRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Errorf("Failed to unmarshal CompletionRequest: %v", err)
	}

	if decoded.Model != req.Model {
		t.Errorf("Decoded Model = %q, want %q", decoded.Model, req.Model)
	}

	if len(decoded.Messages) != len(req.Messages) {
		t.Errorf("Decoded Messages length = %d, want %d", len(decoded.Messages), len(req.Messages))
	}
}

func TestConstants(t *testing.T) {
	if DefaultAPIEndpoint == "" {
		t.Error("DefaultAPIEndpoint should not be empty")
	}

	if DefaultModel == "" {
		t.Error("DefaultModel should not be empty")
	}

	if MaxTokens <= 0 {
		t.Error("MaxTokens should be positive")
	}

	if DefaultTimeout <= 0 {
		t.Error("DefaultTimeout should be positive")
	}
}
