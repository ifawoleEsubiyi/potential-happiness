package webhook

import (
	"testing"
)

func TestNew(t *testing.T) {
	w := New()
	if w == nil {
		t.Fatal("New() returned nil")
	}

	if w.GetURL() != "" {
		t.Errorf("New() URL = %q, want empty", w.GetURL())
	}

	if w.IsEnabled() {
		t.Error("New() IsEnabled = true, want false")
	}

	if w.IsValidated() {
		t.Error("New() IsValidated = true, want false")
	}
}

func TestNewWithURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			url:     "https://api.miningpool.com/webhook",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL with path",
			url:     "https://webhook.example.com/bitcoin/mining/notify",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "HTTP URL (not HTTPS)",
			url:     "http://api.miningpool.com/webhook",
			wantErr: true,
		},
		{
			name:    "localhost URL",
			url:     "https://localhost/webhook",
			wantErr: true,
		},
		{
			name:    "private IP URL",
			url:     "https://192.168.1.1/webhook",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "not-a-valid-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWithURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && w == nil {
				t.Error("NewWithURL() returned nil for valid URL")
			}
			if !tt.wantErr && w.GetURL() != tt.url {
				t.Errorf("NewWithURL() URL = %q, want %q", w.GetURL(), tt.url)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid HTTPS URL",
			url:     "https://api.miningpool.com/webhook",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL with port",
			url:     "https://api.miningpool.com:8443/webhook",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL with query params",
			url:     "https://api.miningpool.com/webhook?token=abc",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
			errMsg:  "webhook URL cannot be empty",
		},
		{
			name:    "HTTP URL not allowed",
			url:     "http://api.miningpool.com/webhook",
			wantErr: true,
			errMsg:  "webhook URL must use HTTPS protocol",
		},
		{
			name:    "FTP URL not allowed",
			url:     "ftp://api.miningpool.com/webhook",
			wantErr: true,
			errMsg:  "webhook URL must use HTTPS protocol",
		},
		{
			name:    "localhost not allowed",
			url:     "https://localhost/webhook",
			wantErr: true,
			errMsg:  "webhook URL cannot point to localhost or private IP ranges",
		},
		{
			name:    "127.0.0.1 not allowed",
			url:     "https://127.0.0.1/webhook",
			wantErr: true,
			errMsg:  "webhook URL cannot point to localhost or private IP ranges",
		},
		{
			name:    "10.x.x.x not allowed",
			url:     "https://10.0.0.1/webhook",
			wantErr: true,
			errMsg:  "webhook URL cannot point to localhost or private IP ranges",
		},
		{
			name:    "192.168.x.x not allowed",
			url:     "https://192.168.1.1/webhook",
			wantErr: true,
			errMsg:  "webhook URL cannot point to localhost or private IP ranges",
		},
		{
			name:    "172.16.x.x not allowed",
			url:     "https://172.16.0.1/webhook",
			wantErr: true,
			errMsg:  "webhook URL cannot point to localhost or private IP ranges",
		},
		{
			name:    "missing host",
			url:     "https:///webhook",
			wantErr: true,
			errMsg:  "webhook URL must have a valid host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("ValidateURL() error = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSetURL(t *testing.T) {
	w := New()

	// Test setting a valid URL
	err := w.SetURL("https://api.miningpool.com/webhook")
	if err != nil {
		t.Errorf("SetURL() error = %v, want nil", err)
	}

	if w.GetURL() != "https://api.miningpool.com/webhook" {
		t.Errorf("SetURL() URL = %q, want %q", w.GetURL(), "https://api.miningpool.com/webhook")
	}

	if !w.IsValidated() {
		t.Error("SetURL() IsValidated = false, want true")
	}

	if !w.IsEnabled() {
		t.Error("SetURL() IsEnabled = false, want true")
	}

	// Test setting an invalid URL
	err = w.SetURL("http://insecure.com/webhook")
	if err == nil {
		t.Error("SetURL() with HTTP URL should return error")
	}
}

func TestEnableDisable(t *testing.T) {
	w := New()

	// Cannot enable without validated URL
	err := w.Enable()
	if err == nil {
		t.Error("Enable() without validated URL should return error")
	}

	// Set valid URL
	err = w.SetURL("https://api.miningpool.com/webhook")
	if err != nil {
		t.Fatalf("SetURL() error = %v", err)
	}

	// Now disable and re-enable should work
	w.Disable()
	if w.IsEnabled() {
		t.Error("Disable() IsEnabled = true, want false")
	}

	err = w.Enable()
	if err != nil {
		t.Errorf("Enable() error = %v, want nil", err)
	}

	if !w.IsEnabled() {
		t.Error("Enable() IsEnabled = false, want true")
	}
}

func TestSubscribeToEvent(t *testing.T) {
	w := New()

	// Subscribe to mining reward event
	err := w.SubscribeToEvent(EventMiningReward)
	if err != nil {
		t.Errorf("SubscribeToEvent() error = %v, want nil", err)
	}

	events := w.GetSubscribedEvents()
	if len(events) != 1 {
		t.Errorf("GetSubscribedEvents() len = %d, want 1", len(events))
	}

	if events[0] != EventMiningReward {
		t.Errorf("GetSubscribedEvents()[0] = %v, want %v", events[0], EventMiningReward)
	}

	// Subscribe to same event again (should not duplicate)
	err = w.SubscribeToEvent(EventMiningReward)
	if err != nil {
		t.Errorf("SubscribeToEvent() duplicate error = %v, want nil", err)
	}

	events = w.GetSubscribedEvents()
	if len(events) != 1 {
		t.Errorf("GetSubscribedEvents() after duplicate len = %d, want 1", len(events))
	}

	// Subscribe to block found event
	err = w.SubscribeToEvent(EventBlockFound)
	if err != nil {
		t.Errorf("SubscribeToEvent() error = %v, want nil", err)
	}

	events = w.GetSubscribedEvents()
	if len(events) != 2 {
		t.Errorf("GetSubscribedEvents() len = %d, want 2", len(events))
	}
}

func TestUnsubscribeFromEvent(t *testing.T) {
	w := New()

	// Subscribe to multiple events
	_ = w.SubscribeToEvent(EventMiningReward)
	_ = w.SubscribeToEvent(EventBlockFound)
	_ = w.SubscribeToEvent(EventPoolPayout)

	events := w.GetSubscribedEvents()
	if len(events) != 3 {
		t.Fatalf("Initial events len = %d, want 3", len(events))
	}

	// Unsubscribe from block found
	w.UnsubscribeFromEvent(EventBlockFound)

	events = w.GetSubscribedEvents()
	if len(events) != 2 {
		t.Errorf("GetSubscribedEvents() after unsubscribe len = %d, want 2", len(events))
	}

	// Verify block found is not in the list
	for _, e := range events {
		if e == EventBlockFound {
			t.Error("EventBlockFound should not be in events after unsubscribe")
		}
	}
}

func TestGetConfig(t *testing.T) {
	w, err := NewWithURL("https://api.miningpool.com/webhook")
	if err != nil {
		t.Fatalf("NewWithURL() error = %v", err)
	}

	config := w.GetConfig()

	if config.URL != "https://api.miningpool.com/webhook" {
		t.Errorf("GetConfig() URL = %q, want %q", config.URL, "https://api.miningpool.com/webhook")
	}

	if !config.Enabled {
		t.Error("GetConfig() Enabled = false, want true")
	}

	if len(config.Events) != 2 {
		t.Errorf("GetConfig() Events len = %d, want 2", len(config.Events))
	}
}

func TestGetConfigEventsCopy(t *testing.T) {
	w, err := NewWithURL("https://api.miningpool.com/webhook")
	if err != nil {
		t.Fatalf("NewWithURL() error = %v", err)
	}

	// Get config and modify the returned Events slice
	config := w.GetConfig()
	originalLen := len(config.Events)
	config.Events = append(config.Events, EventHashrateChange)

	// Get config again and verify internal state was not modified
	config2 := w.GetConfig()
	if len(config2.Events) != originalLen {
		t.Errorf("GetConfig() Events len after external modification = %d, want %d", len(config2.Events), originalLen)
	}
}

func TestEventTypes(t *testing.T) {
	tests := []struct {
		event EventType
		want  string
	}{
		{EventMiningReward, "mining_reward"},
		{EventBlockFound, "block_found"},
		{EventPoolPayout, "pool_payout"},
		{EventHashrateChange, "hashrate_change"},
	}

	for _, tt := range tests {
		t.Run(string(tt.event), func(t *testing.T) {
			if string(tt.event) != tt.want {
				t.Errorf("EventType = %q, want %q", tt.event, tt.want)
			}
		})
	}
}

func TestURLLengthLimit(t *testing.T) {
	// Create a URL that exceeds MaxURLLength
	longPath := make([]byte, MaxURLLength+1)
	for i := range longPath {
		longPath[i] = 'a'
	}

	longURL := "https://example.com/" + string(longPath)
	err := ValidateURL(longURL)
	if err == nil {
		t.Error("ValidateURL() with long URL should return error")
	}

	if err.Error() != "webhook URL exceeds maximum length" {
		t.Errorf("ValidateURL() error = %q, want %q", err.Error(), "webhook URL exceeds maximum length")
	}
}
