// Package farming provides farming functionality for the validator daemon,
// including onboarding and payout management.
package farming

import (
	"errors"
	"sync"

	"github.com/dreadwitdastacc-IFA/validatord/internal/payment"
)

// PayoutSchedule defines when payouts should occur
type PayoutSchedule string

const (
	// PayoutDaily means payouts occur daily (default)
	PayoutDaily PayoutSchedule = "daily"
	// PayoutWeekly means payouts occur weekly
	PayoutWeekly PayoutSchedule = "weekly"
	// PayoutMonthly means payouts occur monthly
	PayoutMonthly PayoutSchedule = "monthly"
)

// FarmingConfig holds the configuration for farming operations
type FarmingConfig struct {
	PayoutSchedule PayoutSchedule
	MinimumPayout  float64
	Enabled        bool
}

// Farmer manages farming operations for validators.
// This type is safe for concurrent use by multiple goroutines.
type Farmer struct {
	mu        sync.RWMutex
	config    FarmingConfig
	onboarded bool
}

// New creates a new farmer instance with default configuration (daily payouts)
func New() *Farmer {
	return &Farmer{
		config: FarmingConfig{
			PayoutSchedule: PayoutDaily,
			MinimumPayout:  0.0,
			Enabled:        false,
		},
		onboarded: false,
	}
}

// NewWithConfig creates a new farmer instance with custom configuration
func NewWithConfig(config FarmingConfig) *Farmer {
	return &Farmer{
		config:    config,
		onboarded: false,
	}
}

// Onboard performs the onboarding process for a new farmer
func (f *Farmer) Onboard(paystring string) error {
	if paystring == "" {
		return errors.New("paystring is required for onboarding")
	}

	// Validate the paystring format before acquiring lock
	if err := payment.ValidatePaystring(paystring); err != nil {
		return err
	}

	// Acquire lock and check if already onboarded
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.onboarded {
		return errors.New("farmer is already onboarded")
	}

	f.onboarded = true
	f.config.Enabled = true
	return nil
}

// IsOnboarded returns whether the farmer has been onboarded
func (f *Farmer) IsOnboarded() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.onboarded
}

// GetConfig returns the current farming configuration
func (f *Farmer) GetConfig() FarmingConfig {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.config
}

// UpdateConfig updates the farming configuration
func (f *Farmer) UpdateConfig(config FarmingConfig) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.onboarded {
		return errors.New("farmer must be onboarded before updating config")
	}
	f.config = config
	return nil
}

// IsEnabled returns whether farming is enabled
func (f *Farmer) IsEnabled() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.config.Enabled
}

// GetPayoutSchedule returns the current payout schedule
func (f *Farmer) GetPayoutSchedule() PayoutSchedule {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.config.PayoutSchedule
}
