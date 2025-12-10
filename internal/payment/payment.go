// Package payment provides payment and paystring functionality for the validator daemon.
package payment

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// redactedPrefixLength is the number of characters to show before redacting the rest
	redactedPrefixLength = 3

	// paystringSeparator is the separator used in paystring format (user$domain.com)
	paystringSeparator = "$"

	// PayoutStatusCompleted indicates a payout was successfully completed
	PayoutStatusCompleted = "completed"
	// PayoutStatusPending indicates a payout is being processed
	PayoutStatusPending = "pending"
	// PayoutStatusFailed indicates a payout has failed
	PayoutStatusFailed = "failed"

	// MaxPayoutAmount is the maximum allowed payout amount in smallest currency unit (e.g., cents)
	MaxPayoutAmount = 100000000 // 1,000,000.00 in major currency units

	// MaxReferenceLength is the maximum allowed length for a payout reference
	MaxReferenceLength = 256
)

// splitPaystring splits a paystring into user and domain parts.
// Returns the user, domain, and true if the paystring has exactly one separator
// and both parts are non-empty. Returns empty strings and false if the format is invalid.
func splitPaystring(paystring string) (user, domain string, ok bool) {
	paystringParts := strings.Split(paystring, paystringSeparator)
	if len(paystringParts) != 2 || paystringParts[0] == "" || paystringParts[1] == "" {
		return "", "", false
	}
	return paystringParts[0], paystringParts[1], true
}

// Payment manages payment information for validators.
type Payment struct {
	paystring string
}

// PayoutRequest represents a request to execute a payout.
type PayoutRequest struct {
	Recipient string // Recipient paystring in format user$domain.com
	Amount    int64  // Amount in smallest currency unit (e.g., cents for USD, pence for GBP)
	Currency  string // ISO 4217 currency code (e.g., USD, EUR, GBP)
	Reference string // Optional reference or description for the payout. If provided, must be a UTF-8 string up to 256 characters. Control characters are not allowed.
}

// PayoutResult represents the result of a payout execution.
type PayoutResult struct {
	TransactionID string
	Status        string
	StatusCode    int
	Timestamp     time.Time
}

// New creates a new payment instance with the provided paystring.
func New(paystring string) (*Payment, error) {
	if err := ValidatePaystring(paystring); err != nil {
		return nil, err
	}
	return &Payment{
		paystring: paystring,
	}, nil
}

// GetPaystring returns the configured paystring.
func (p *Payment) GetPaystring() string {
	return p.paystring
}

// GetRedactedPaystring returns the paystring with the user portion partially redacted.
// For example, "ifawoleesubiyi$paystring.crypto.com" becomes "ifa***$paystring.crypto.com"
func (p *Payment) GetRedactedPaystring() string {
	return RedactPaystring(p.paystring)
}

// RedactPaystring redacts a paystring by masking most of the user portion.
// It keeps up to the first 3 characters visible and adds asterisks.
// For example, "ifawoleesubiyi$paystring.crypto.com" becomes "ifa***$paystring.crypto.com"
// For shorter names, "bob$example.com" becomes "bob***$example.com"
func RedactPaystring(paystring string) string {
	user, domain, ok := splitPaystring(paystring)
	if !ok {
		// If invalid format, redact entire string
		return "***"
	}

	// Show up to first redactedPrefixLength characters, add asterisks
	var redactedUser string
	if len(user) <= redactedPrefixLength {
		redactedUser = user + "***"
	} else {
		redactedUser = user[:redactedPrefixLength] + "***"
	}

	return redactedUser + paystringSeparator + domain
}

// ValidatePaystring validates that a paystring is in the correct format.
// Paystrings should be in the format: user$domain.com
func ValidatePaystring(paystring string) error {
	if paystring == "" {
		return errors.New("paystring cannot be empty")
	}

	// Check for correct number of separators first
	paystringParts := strings.Split(paystring, paystringSeparator)
	if len(paystringParts) != 2 {
		return errors.New("paystring must contain exactly one '" + paystringSeparator + "' separator")
	}

	// Check for empty parts with specific error messages
	if paystringParts[0] == "" {
		return errors.New("paystring user portion cannot be empty")
	}

	if paystringParts[1] == "" {
		return errors.New("paystring domain portion cannot be empty")
	}

	if !strings.Contains(paystringParts[1], ".") {
		return errors.New("paystring domain must contain at least one '.'")
	}

	return nil
}

// ValidatePayoutRequest validates a payout request.
func ValidatePayoutRequest(req *PayoutRequest) error {
	if req == nil {
		return errors.New("payout request cannot be nil")
	}

	if err := ValidatePaystring(req.Recipient); err != nil {
		return fmt.Errorf("invalid recipient paystring: %w", err)
	}

	if req.Amount <= 0 {
		return errors.New("payout amount must be greater than zero")
	}

	if req.Amount > MaxPayoutAmount {
		return fmt.Errorf("payout amount exceeds maximum allowed: %d", MaxPayoutAmount)
	}

	if req.Currency == "" {
		return errors.New("payout currency cannot be empty")
	}

	if !isValidCurrency(req.Currency) {
		return errors.New("payout currency must be a valid 3-letter ISO 4217 code")
	}

	// Validate Reference field if provided
	if req.Reference != "" {
		if len(req.Reference) > MaxReferenceLength {
			return errors.New("payout reference exceeds maximum length of 256 characters")
		}
		for _, r := range req.Reference {
			if r < 32 || r == 127 {
				return errors.New("payout reference contains control characters")
			}
		}
	}

	return nil
}

// isValidCurrency validates that a currency code matches ISO 4217 format (3 uppercase letters).
// Note: This does not verify the code is a registered ISO 4217 currency.
func isValidCurrency(currency string) bool {
	// Basic length check for ISO 4217 (3-letter codes)
	if len(currency) != 3 {
		return false
	}
	// Currency codes should be uppercase alphabetic
	for _, char := range currency {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}

// generateTransactionID creates a unique transaction ID using cryptographically secure random bytes.
func generateTransactionID() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return fmt.Sprintf("tx_%s", hex.EncodeToString(randomBytes)), nil
}

// ExecutePayout executes a payout to the specified recipient.
func (p *Payment) ExecutePayout(req *PayoutRequest) (*PayoutResult, error) {
	if err := ValidatePayoutRequest(req); err != nil {
		return nil, fmt.Errorf("invalid payout request: %w", err)
	}

	// Generate a unique transaction ID
	txID, err := generateTransactionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ID: %w", err)
	}

	// In a real implementation, this would interact with a payment provider
	// For now, we simulate a successful payout
	result := &PayoutResult{
		TransactionID: txID,
		Status:        PayoutStatusCompleted,
		StatusCode:    http.StatusOK,
		Timestamp:     time.Now(),
	}

	return result, nil
}

// ProcessPayout processes a payout request and returns the result.
// This is a convenience wrapper around ExecutePayout.
// amount is in smallest currency unit (e.g., cents for USD).
func (p *Payment) ProcessPayout(recipient string, amount int64, currency string) (*PayoutResult, error) {
	req := &PayoutRequest{
		Recipient: recipient,
		Amount:    amount,
		Currency:  currency,
		Reference: fmt.Sprintf("payout_%d", time.Now().Unix()),
	}

	return p.ExecutePayout(req)
}
