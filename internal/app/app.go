// Package app provides the shared initialization and running logic for the validatord daemon.
package app

import (
	"fmt"

	"github.com/dreadwitdastacc-IFA/validatord/internal/aggregator"
	"github.com/dreadwitdastacc-IFA/validatord/internal/attest"
	"github.com/dreadwitdastacc-IFA/validatord/internal/bls"
	"github.com/dreadwitdastacc-IFA/validatord/internal/farming"
	"github.com/dreadwitdastacc-IFA/validatord/internal/keystore"
	"github.com/dreadwitdastacc-IFA/validatord/internal/llm"
	"github.com/dreadwitdastacc-IFA/validatord/internal/milestone"
	"github.com/dreadwitdastacc-IFA/validatord/internal/models"
	"github.com/dreadwitdastacc-IFA/validatord/internal/payment"
	"github.com/dreadwitdastacc-IFA/validatord/internal/watcher"
	"github.com/dreadwitdastacc-IFA/validatord/internal/webhook"
)

// DefaultPaystring is the default paystring for the validatord daemon.
const DefaultPaystring = "ifawoleesubiyi$paystring.crypto.com"

// App represents the validatord application with all its components.
type App struct {
	Payment    *payment.Payment
	Farmer     *farming.Farmer
	Attest     *attest.Attest
	Aggregator *aggregator.Aggregator
	BLS        *bls.BLS
	Keystore   *keystore.Keystore
	Watcher    *watcher.Watcher
	Webhook    *webhook.Webhook
	Milestone  *milestone.Maker
	Models     *models.Models
	LLM        *llm.Client
}

// New creates and initializes a new validatord application with all components.
func New(paystring string) (*App, error) {
	// Initialize payment with the configured paystring
	pay, err := payment.New(paystring)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize payment: %w", err)
	}

	// Initialize farming with default daily payout template
	farmer := farming.New()
	if err := farmer.Onboard(pay.GetPaystring()); err != nil {
		return nil, fmt.Errorf("failed to onboard farmer: %w", err)
	}

	return &App{
		Payment:    pay,
		Farmer:     farmer,
		Attest:     attest.New(),
		Aggregator: aggregator.New(),
		BLS:        bls.New(),
		Keystore:   keystore.New(),
		Watcher:    watcher.New(),
		Webhook:    webhook.New(),
		Milestone:  milestone.New(),
		Models:     models.New(),
		LLM:        llm.New(),
	}, nil
}

// PrintStatus prints the current status of the application.
func (a *App) PrintStatus() {
	fmt.Println("Validatord initialized successfully")
	fmt.Printf("Payment address (redacted): %s\n", a.Payment.GetRedactedPaystring())
	fmt.Printf("Farming enabled: %v\n", a.Farmer.IsEnabled())
	fmt.Printf("Payout schedule: %s\n", a.Farmer.GetPayoutSchedule())
	fmt.Printf("Webhook enabled: %v\n", a.Webhook.IsEnabled())
	fmt.Printf("Milestones tracked: %d\n", a.Milestone.Count())
	fmt.Printf("GitHub Models enabled: %v\n", a.Models.HasToken())
	fmt.Printf("Default model: %s\n", a.Models.GetDefaultModel())
	fmt.Printf("LLM API endpoint: %s\n", a.LLM.GetAPIEndpoint())
	fmt.Printf("LLM model: %s\n", a.LLM.GetModel())
	fmt.Printf("LLM token configured: %v\n", a.LLM.HasToken())
}
