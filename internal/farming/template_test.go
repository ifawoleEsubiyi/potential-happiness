package farming

import (
	"testing"
)

func TestDefaultTemplates(t *testing.T) {
	templates := DefaultTemplates()

	if len(templates) == 0 {
		t.Fatal("DefaultTemplates() returned empty slice")
	}

	// Check that we have the expected templates
	expectedNames := []string{"Daily Payouts", "Weekly Payouts", "Monthly Payouts"}
	if len(templates) != len(expectedNames) {
		t.Errorf("Expected %d templates, got %d", len(expectedNames), len(templates))
	}

	// Verify each template has required fields
	for _, tmpl := range templates {
		if tmpl.Name == "" {
			t.Error("Template has empty Name")
		}
		if tmpl.Description == "" {
			t.Error("Template has empty Description")
		}
		if tmpl.PayoutSchedule == "" {
			t.Error("Template has empty PayoutSchedule")
		}
	}
}

func TestGetTemplateByName(t *testing.T) {
	tests := []struct {
		name      string
		tmplName  string
		wantFound bool
		wantSched PayoutSchedule
	}{
		{
			name:      "daily payouts template",
			tmplName:  "Daily Payouts",
			wantFound: true,
			wantSched: PayoutDaily,
		},
		{
			name:      "weekly payouts template",
			tmplName:  "Weekly Payouts",
			wantFound: true,
			wantSched: PayoutWeekly,
		},
		{
			name:      "monthly payouts template",
			tmplName:  "Monthly Payouts",
			wantFound: true,
			wantSched: PayoutMonthly,
		},
		{
			name:      "non-existent template",
			tmplName:  "Non-Existent",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, found := GetTemplateByName(tt.tmplName)

			if found != tt.wantFound {
				t.Errorf("GetTemplateByName() found = %v, want %v", found, tt.wantFound)
				return
			}

			if tt.wantFound {
				if tmpl == nil {
					t.Fatal("GetTemplateByName() returned nil template when found was true")
				}

				if tmpl.Name != tt.tmplName {
					t.Errorf("Template name = %s, want %s", tmpl.Name, tt.tmplName)
				}

				if tmpl.PayoutSchedule != tt.wantSched {
					t.Errorf("Template schedule = %s, want %s", tmpl.PayoutSchedule, tt.wantSched)
				}
			} else {
				if tmpl != nil {
					t.Error("GetTemplateByName() should return nil template when not found")
				}
			}
		})
	}
}

func TestApplyTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template Template
	}{
		{
			name: "daily payout template",
			template: Template{
				Name:           "Daily Payouts",
				Description:    "Default template with daily payouts",
				PayoutSchedule: PayoutDaily,
				MinimumPayout:  0.0,
			},
		},
		{
			name: "weekly payout template",
			template: Template{
				Name:           "Weekly Payouts",
				Description:    "Weekly payout schedule",
				PayoutSchedule: PayoutWeekly,
				MinimumPayout:  5.0,
			},
		},
		{
			name: "monthly payout template",
			template: Template{
				Name:           "Monthly Payouts",
				Description:    "Monthly payout schedule",
				PayoutSchedule: PayoutMonthly,
				MinimumPayout:  10.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ApplyTemplate(tt.template)

			if config.PayoutSchedule != tt.template.PayoutSchedule {
				t.Errorf("Config schedule = %s, want %s", config.PayoutSchedule, tt.template.PayoutSchedule)
			}

			if config.MinimumPayout != tt.template.MinimumPayout {
				t.Errorf("Config minimum payout = %f, want %f", config.MinimumPayout, tt.template.MinimumPayout)
			}

			// Applied template should always start with Enabled=false
			if config.Enabled {
				t.Error("Applied template config should have Enabled=false")
			}
		})
	}
}

func TestDefaultTemplate_DailyPayouts(t *testing.T) {
	// Test that the default "Daily Payouts" template exists and has correct values
	tmpl, found := GetTemplateByName("Daily Payouts")
	if !found {
		t.Fatal("Default 'Daily Payouts' template not found")
	}

	if tmpl.PayoutSchedule != PayoutDaily {
		t.Errorf("Daily Payouts template should have PayoutDaily schedule, got %s", tmpl.PayoutSchedule)
	}

	if tmpl.MinimumPayout != 0.2 {
		t.Errorf("Daily Payouts template should have MinimumPayout 0.2, got %f", tmpl.MinimumPayout)
	}

	// Apply the template and verify it creates a valid config
	config := ApplyTemplate(*tmpl)
	if config.PayoutSchedule != PayoutDaily {
		t.Errorf("Applied config should have PayoutDaily schedule, got %s", config.PayoutSchedule)
	}
}
