package models

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.HasToken() {
		t.Error("New() should not have a token configured")
	}
	if m.GetEndpoint() != DefaultAPIEndpoint {
		t.Errorf("GetEndpoint() = %v, want %v", m.GetEndpoint(), DefaultAPIEndpoint)
	}
	if m.GetDefaultModel() != DefaultModel {
		t.Errorf("GetDefaultModel() = %v, want %v", m.GetDefaultModel(), DefaultModel)
	}
}

func TestNewWithToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid token",
			token:   "ghp_1234567890abcdef",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			errMsg:  "GitHub token cannot be empty",
		},
		{
			name:    "short token",
			token:   "short",
			wantErr: true,
			errMsg:  "GitHub token appears to be invalid (too short)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewWithToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("NewWithToken() error = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr {
				if m == nil {
					t.Fatal("NewWithToken() returned nil with no error")
				}
				if !m.HasToken() {
					t.Error("NewWithToken() should have token configured")
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid token",
			token:   "ghp_1234567890abcdefghijklmnopqrstuvwxyz",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			errMsg:  "GitHub token cannot be empty",
		},
		{
			name:    "token too short",
			token:   "abc123",
			wantErr: true,
			errMsg:  "GitHub token appears to be invalid (too short)",
		},
		{
			name:    "token with whitespace",
			token:   "  ghp_1234567890abcdef  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValidateToken() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSetToken(t *testing.T) {
	m := New()

	// Test setting valid token
	token := "ghp_validtoken123"
	err := m.SetToken(token)
	if err != nil {
		t.Errorf("SetToken() unexpected error: %v", err)
	}
	if !m.HasToken() {
		t.Error("SetToken() should set HasToken to true")
	}

	// Test setting empty token
	err = m.SetToken("")
	if err == nil {
		t.Error("SetToken() should return error for empty token")
	}
}

func TestSetDefaultModel(t *testing.T) {
	m := New()

	tests := []struct {
		name    string
		model   string
		wantErr bool
	}{
		{
			name:    "valid model",
			model:   "openai/gpt-4",
			wantErr: false,
		},
		{
			name:    "empty model",
			model:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.SetDefaultModel(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultModel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && m.GetDefaultModel() != tt.model {
				t.Errorf("GetDefaultModel() = %v, want %v", m.GetDefaultModel(), tt.model)
			}
		})
	}
}

func TestValidateChatRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *ChatRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "chat request cannot be nil",
		},
		{
			name: "empty model",
			req: &ChatRequest{
				Model: "",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
			errMsg:  "model must be specified",
		},
		{
			name: "empty messages",
			req: &ChatRequest{
				Model:    "openai/gpt-4",
				Messages: []Message{},
			},
			wantErr: true,
			errMsg:  "messages cannot be empty",
		},
		{
			name: "message with empty role",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "", Content: "Hello"},
				},
			},
			wantErr: true,
			errMsg:  "message 0: role cannot be empty",
		},
		{
			name: "message with empty content",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: ""},
				},
			},
			wantErr: true,
			errMsg:  "message 0: content cannot be empty",
		},
		{
			name: "message with invalid role",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "invalid", Content: "Hello"},
				},
			},
			wantErr: true,
			errMsg:  "message 0: invalid role 'invalid' (must be 'user', 'assistant', or 'system')",
		},
		{
			name: "message content too long",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: strings.Repeat("a", MaxPromptLength+1)},
				},
			},
			wantErr: true,
		},
		{
			name: "temperature too low",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: -0.1,
			},
			wantErr: true,
			errMsg:  "temperature must be between 0.0 and 2.0",
		},
		{
			name: "temperature too high",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				Temperature: 2.1,
			},
			wantErr: true,
			errMsg:  "temperature must be between 0.0 and 2.0",
		},
		{
			name: "negative max tokens",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens: -1,
			},
			wantErr: true,
			errMsg:  "max_tokens must be non-negative",
		},
		{
			name: "valid system message",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "system", Content: "You are a helpful assistant"},
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid assistant message",
			req: &ChatRequest{
				Model: "openai/gpt-4",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
					{Role: "assistant", Content: "Hi there!"},
					{Role: "user", Content: "How are you?"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChatRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChatRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("ValidateChatRequest() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCallModel_NoToken(t *testing.T) {
	m := New()
	req := &ChatRequest{
		Model: "openai/gpt-4",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	_, err := m.CallModel(req)
	if err == nil {
		t.Error("CallModel() should return error when token not configured")
	}
	if err != nil && !strings.Contains(err.Error(), "token not configured") {
		t.Errorf("CallModel() error = %v, want error containing 'token not configured'", err)
	}
}

func TestCallModel_InvalidRequest(t *testing.T) {
	m := New()
	m.SetToken("ghp_validtoken123")

	req := &ChatRequest{
		Model:    "",
		Messages: []Message{},
	}

	_, err := m.CallModel(req)
	if err == nil {
		t.Error("CallModel() should return error for invalid request")
	}
}

func TestChat_EmptyPrompt(t *testing.T) {
	m := New()
	m.SetToken("ghp_validtoken123")

	_, err := m.Chat("")
	if err == nil {
		t.Error("Chat() should return error for empty prompt")
	}
	if err != nil && !strings.Contains(err.Error(), "prompt cannot be empty") {
		t.Errorf("Chat() error = %v, want error containing 'prompt cannot be empty'", err)
	}
}

func TestChatWithModel_EmptyInputs(t *testing.T) {
	m := New()
	m.SetToken("ghp_validtoken123")

	tests := []struct {
		name    string
		prompt  string
		model   string
		wantErr string
	}{
		{
			name:    "empty prompt",
			prompt:  "",
			model:   "openai/gpt-4",
			wantErr: "prompt cannot be empty",
		},
		{
			name:    "empty model",
			prompt:  "Hello",
			model:   "",
			wantErr: "model cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.ChatWithModel(tt.prompt, tt.model)
			if err == nil {
				t.Error("ChatWithModel() should return error")
			}
			if err != nil && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ChatWithModel() error = %v, want error containing '%v'", err, tt.wantErr)
			}
		})
	}
}

func TestGetters(t *testing.T) {
	m := New()

	// Test default values
	if m.GetEndpoint() != DefaultAPIEndpoint {
		t.Errorf("GetEndpoint() = %v, want %v", m.GetEndpoint(), DefaultAPIEndpoint)
	}
	if m.GetDefaultModel() != DefaultModel {
		t.Errorf("GetDefaultModel() = %v, want %v", m.GetDefaultModel(), DefaultModel)
	}
	if m.HasToken() {
		t.Error("HasToken() should return false for new instance")
	}

	// Test after setting token
	m.SetToken("ghp_validtoken123")
	if !m.HasToken() {
		t.Error("HasToken() should return true after SetToken")
	}

	// Test after setting model
	newModel := "openai/gpt-4"
	m.SetDefaultModel(newModel)
	if m.GetDefaultModel() != newModel {
		t.Errorf("GetDefaultModel() = %v, want %v", m.GetDefaultModel(), newModel)
	}
}

func TestMessageStructure(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	if msg.Role != "user" {
		t.Errorf("Message.Role = %v, want user", msg.Role)
	}
	if msg.Content != "Hello, world!" {
		t.Errorf("Message.Content = %v, want 'Hello, world!'", msg.Content)
	}
}
