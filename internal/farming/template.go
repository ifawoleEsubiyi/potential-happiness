// Package farming provides farming functionality for the validator daemon,
// including onboarding and payout management.
package farming

// Template represents a farming configuration template
type Template struct {
	Name           string
	Description    string
	PayoutSchedule PayoutSchedule
	MinimumPayout  float64
}

// defaultTemplates is the cached list of default farming templates.
// This variable should not be modified after initialization. All functions return copies to prevent external modifications.
// Note: Go does not enforce read-only variables; this is a convention, not a language-enforced guarantee.
var defaultTemplates = []Template{
	{
		Name:           "Daily Payouts",
		Description:    "Default template with daily payouts, suitable for most validators",
		PayoutSchedule: PayoutDaily,
		MinimumPayout:  0.2,
	},
	{
		Name:           "Weekly Payouts",
		Description:    "Weekly payout schedule for reduced transaction overhead",
		PayoutSchedule: PayoutWeekly,
		MinimumPayout:  0.2,
	},
	{
		Name:           "Monthly Payouts",
		Description:    "Monthly payout schedule for maximum efficiency",
		PayoutSchedule: PayoutMonthly,
		MinimumPayout:  0.2,
	},
}

// DefaultTemplates returns the default farming templates
func DefaultTemplates() []Template {
	// Return a copy to prevent external modifications
	templatesCopy := make([]Template, len(defaultTemplates))
	copy(templatesCopy, defaultTemplates)
	return templatesCopy
}

// GetTemplateByName returns a template by name
func GetTemplateByName(name string) (*Template, bool) {
	for _, template := range defaultTemplates {
		if template.Name == name {
			// Return a copy to prevent external modifications
			templateCopy := template
			return &templateCopy, true
		}
	}
	return nil, false
}

// ApplyTemplate applies a template to create a farming configuration
func ApplyTemplate(template Template) FarmingConfig {
	return FarmingConfig{
		PayoutSchedule: template.PayoutSchedule,
		MinimumPayout:  template.MinimumPayout,
		Enabled:        false, // Must be explicitly enabled after onboarding
	}
}
