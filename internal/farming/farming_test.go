package farming

import (
	"testing"
)

func TestNew(t *testing.T) {
	f := New()
	if f == nil {
		t.Fatal("New() returned nil")
	}

	if f.IsOnboarded() {
		t.Error("New farmer should not be onboarded")
	}

	if f.IsEnabled() {
		t.Error("New farmer should not be enabled")
	}

	if f.GetPayoutSchedule() != PayoutDaily {
		t.Errorf("New farmer should have daily payout schedule, got %s", f.GetPayoutSchedule())
	}
}

func TestNewWithConfig(t *testing.T) {
	config := FarmingConfig{
		PayoutSchedule: PayoutWeekly,
		MinimumPayout:  10.0,
		Enabled:        false,
	}

	f := NewWithConfig(config)
	if f == nil {
		t.Fatal("NewWithConfig() returned nil")
	}

	if f.GetPayoutSchedule() != PayoutWeekly {
		t.Errorf("Expected weekly payout schedule, got %s", f.GetPayoutSchedule())
	}

	gotConfig := f.GetConfig()
	if gotConfig.MinimumPayout != 10.0 {
		t.Errorf("Expected minimum payout 10.0, got %f", gotConfig.MinimumPayout)
	}
}

func TestOnboard(t *testing.T) {
	tests := []struct {
		name       string
		paystring  string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "valid onboarding",
			paystring: "user$domain.com",
			wantErr:   false,
		},
		{
			name:       "empty paystring",
			paystring:  "",
			wantErr:    true,
			errMessage: "paystring is required for onboarding",
		},
		{
			name:       "invalid paystring format - no dollar sign",
			paystring:  "userdomain.com",
			wantErr:    true,
			errMessage: "paystring must contain exactly one '$' separator",
		},
		{
			name:       "invalid paystring format - no domain dot",
			paystring:  "user$domain",
			wantErr:    true,
			errMessage: "paystring domain must contain at least one '.'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New()
			err := f.Onboard(tt.paystring)

			if (err != nil) != tt.wantErr {
				t.Errorf("Onboard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("Onboard() error = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr {
				if !f.IsOnboarded() {
					t.Error("Farmer should be onboarded after successful onboarding")
				}
				if !f.IsEnabled() {
					t.Error("Farmer should be enabled after successful onboarding")
				}
			}
		})
	}
}

func TestOnboard_AlreadyOnboarded(t *testing.T) {
	f := New()

	// First onboarding should succeed
	err := f.Onboard("user$domain.com")
	if err != nil {
		t.Fatalf("First onboarding failed: %v", err)
	}

	// Second onboarding should fail
	err = f.Onboard("user$domain.com")
	if err == nil {
		t.Error("Second onboarding should fail")
	}

	expectedErr := "farmer is already onboarded"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestIsOnboarded(t *testing.T) {
	f := New()

	if f.IsOnboarded() {
		t.Error("New farmer should not be onboarded")
	}

	err := f.Onboard("user$domain.com")
	if err != nil {
		t.Fatalf("Onboard() failed: %v", err)
	}

	if !f.IsOnboarded() {
		t.Error("Farmer should be onboarded after calling Onboard()")
	}
}

func TestUpdateConfig(t *testing.T) {
	f := New()

	// Try to update config before onboarding
	newConfig := FarmingConfig{
		PayoutSchedule: PayoutMonthly,
		MinimumPayout:  5.0,
		Enabled:        true,
	}

	err := f.UpdateConfig(newConfig)
	if err == nil {
		t.Error("UpdateConfig should fail before onboarding")
	}

	// Onboard the farmer
	err = f.Onboard("user$domain.com")
	if err != nil {
		t.Fatalf("Onboard() failed: %v", err)
	}

	// Now update should succeed
	err = f.UpdateConfig(newConfig)
	if err != nil {
		t.Errorf("UpdateConfig() failed after onboarding: %v", err)
	}

	gotConfig := f.GetConfig()
	if gotConfig.PayoutSchedule != PayoutMonthly {
		t.Errorf("Expected monthly payout schedule, got %s", gotConfig.PayoutSchedule)
	}

	if gotConfig.MinimumPayout != 5.0 {
		t.Errorf("Expected minimum payout 5.0, got %f", gotConfig.MinimumPayout)
	}
}

func TestGetPayoutSchedule(t *testing.T) {
	tests := []struct {
		name     string
		schedule PayoutSchedule
	}{
		{
			name:     "daily schedule",
			schedule: PayoutDaily,
		},
		{
			name:     "weekly schedule",
			schedule: PayoutWeekly,
		},
		{
			name:     "monthly schedule",
			schedule: PayoutMonthly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FarmingConfig{
				PayoutSchedule: tt.schedule,
				MinimumPayout:  0.0,
				Enabled:        false,
			}

			f := NewWithConfig(config)
			got := f.GetPayoutSchedule()

			if got != tt.schedule {
				t.Errorf("GetPayoutSchedule() = %s, want %s", got, tt.schedule)
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	f := New()

	if f.IsEnabled() {
		t.Error("New farmer should not be enabled")
	}

	err := f.Onboard("user$domain.com")
	if err != nil {
		t.Fatalf("Onboard() failed: %v", err)
	}

	if !f.IsEnabled() {
		t.Error("Farmer should be enabled after onboarding")
	}
}
