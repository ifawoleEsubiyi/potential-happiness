// Package webhook provides webhook handling functionality for the validator daemon,
// including URL validation and event handling for Bitcoin mining notifications.
package webhook

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"sync"
)

// EventType defines the type of webhook event
type EventType string

const (
	// EventMiningReward indicates a mining reward notification
	EventMiningReward EventType = "mining_reward"
	// EventBlockFound indicates a block was found
	EventBlockFound EventType = "block_found"
	// EventPoolPayout indicates a pool payout event
	EventPoolPayout EventType = "pool_payout"
	// EventHashrateChange indicates a hashrate change notification
	EventHashrateChange EventType = "hashrate_change"

	// MaxURLLength is the maximum allowed length for a webhook URL
	MaxURLLength = 2048
)

// WebhookConfig holds the configuration for webhook operations
type WebhookConfig struct {
	URL     string
	Enabled bool
	Events  []EventType
}

// Webhook manages webhook operations for Bitcoin mining notifications.
// This type is safe for concurrent use by multiple goroutines.
type Webhook struct {
	mu        sync.RWMutex
	config    WebhookConfig
	validated bool
}

// New creates a new webhook instance with default configuration
func New() *Webhook {
	return &Webhook{
		config: WebhookConfig{
			URL:     "",
			Enabled: false,
			Events:  []EventType{},
		},
		validated: false,
	}
}

// NewWithURL creates a new webhook instance with a specified URL
func NewWithURL(webhookURL string) (*Webhook, error) {
	if err := ValidateURL(webhookURL); err != nil {
		return nil, err
	}

	return &Webhook{
		config: WebhookConfig{
			URL:     webhookURL,
			Enabled: true,
			Events:  []EventType{EventMiningReward, EventPoolPayout},
		},
		validated: true,
	}, nil
}

// ValidateURL validates that a webhook URL is properly formatted and secure.
// It checks for HTTPS protocol, valid URL structure, and length limits.
func ValidateURL(webhookURL string) error {
	if webhookURL == "" {
		return errors.New("webhook URL cannot be empty")
	}

	if len(webhookURL) > MaxURLLength {
		return errors.New("webhook URL exceeds maximum length")
	}

	parsed, err := url.Parse(webhookURL)
	if err != nil {
		return errors.New("invalid webhook URL format")
	}

	// Require HTTPS for security
	if parsed.Scheme != "https" {
		return errors.New("webhook URL must use HTTPS protocol")
	}

	// Ensure host is present
	if parsed.Host == "" {
		return errors.New("webhook URL must have a valid host")
	}

	// Check for localhost or internal IPs (basic SSRF protection)
	hostname := parsed.Hostname()
	if strings.EqualFold(hostname, "localhost") {
		return errors.New("webhook URL cannot point to localhost or private IP ranges")
	}

	// Check if the host is an IP address and validate it's not private
	if ip := net.ParseIP(hostname); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return errors.New("webhook URL cannot point to localhost or private IP ranges")
		}
	}

	return nil
}

// SetURL sets and validates a new webhook URL
func (w *Webhook) SetURL(webhookURL string) error {
	if err := ValidateURL(webhookURL); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	w.config.URL = webhookURL
	w.validated = true
	w.config.Enabled = true
	return nil
}

// GetURL returns the configured webhook URL
func (w *Webhook) GetURL() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.config.URL
}

// IsValidated returns whether the webhook URL has been validated
func (w *Webhook) IsValidated() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.validated
}

// IsEnabled returns whether the webhook is enabled
func (w *Webhook) IsEnabled() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.config.Enabled
}

// Enable enables the webhook
func (w *Webhook) Enable() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.validated {
		return errors.New("cannot enable webhook without a validated URL")
	}
	w.config.Enabled = true
	return nil
}

// Disable disables the webhook
func (w *Webhook) Disable() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.config.Enabled = false
}

// GetConfig returns the current webhook configuration
func (w *Webhook) GetConfig() WebhookConfig {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// Create a copy of Events slice to prevent external modification
	events := make([]EventType, len(w.config.Events))
	copy(events, w.config.Events)

	return WebhookConfig{
		URL:     w.config.URL,
		Enabled: w.config.Enabled,
		Events:  events,
	}
}

// SubscribeToEvent adds an event type to the subscription list
func (w *Webhook) SubscribeToEvent(event EventType) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if already subscribed
	for _, e := range w.config.Events {
		if e == event {
			return nil // Already subscribed
		}
	}

	w.config.Events = append(w.config.Events, event)
	return nil
}

// UnsubscribeFromEvent removes an event type from the subscription list
func (w *Webhook) UnsubscribeFromEvent(event EventType) {
	w.mu.Lock()
	defer w.mu.Unlock()

	events := make([]EventType, 0, len(w.config.Events))
	for _, e := range w.config.Events {
		if e != event {
			events = append(events, e)
		}
	}
	w.config.Events = events
}

// GetSubscribedEvents returns the list of subscribed event types
func (w *Webhook) GetSubscribedEvents() []EventType {
	w.mu.RLock()
	defer w.mu.RUnlock()

	events := make([]EventType, len(w.config.Events))
	copy(events, w.config.Events)
	return events
}
