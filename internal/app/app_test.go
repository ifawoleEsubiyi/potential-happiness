package app

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		paystring string
		wantErr   bool
	}{
		{
			name:      "valid paystring",
			paystring: "user$domain.com",
			wantErr:   false,
		},
		{
			name:      "default paystring",
			paystring: DefaultPaystring,
			wantErr:   false,
		},
		{
			name:      "empty paystring",
			paystring: "",
			wantErr:   true,
		},
		{
			name:      "invalid paystring format",
			paystring: "invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(tt.paystring)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if app == nil {
					t.Error("New() returned nil App for valid paystring")
					return
				}
				if app.Payment == nil {
					t.Error("New() returned App with nil Payment")
				}
				if app.Farmer == nil {
					t.Error("New() returned App with nil Farmer")
				}
				if app.Attest == nil {
					t.Error("New() returned App with nil Attest")
				}
				if app.Aggregator == nil {
					t.Error("New() returned App with nil Aggregator")
				}
				if app.BLS == nil {
					t.Error("New() returned App with nil BLS")
				}
				if app.Keystore == nil {
					t.Error("New() returned App with nil Keystore")
				}
				if app.Watcher == nil {
					t.Error("New() returned App with nil Watcher")
				}
				if app.Models == nil {
					t.Error("New() returned App with nil Models")
				}
				if !app.Farmer.IsEnabled() {
					t.Error("Farmer should be enabled after New()")
				}
				if !app.Farmer.IsOnboarded() {
					t.Error("Farmer should be onboarded after New()")
				}
			}
		})
	}
}

func TestDefaultPaystring(t *testing.T) {
	if DefaultPaystring == "" {
		t.Error("DefaultPaystring should not be empty")
	}
}
