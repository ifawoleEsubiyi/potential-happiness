package payment

import (
	"net/http"
	"strings"
	"testing"
)

func TestSplitPaystring(t *testing.T) {
	tests := []struct {
		name       string
		paystring  string
		wantUser   string
		wantDomain string
		wantOk     bool
	}{
		{
			name:       "valid paystring",
			paystring:  "user$domain.com",
			wantUser:   "user",
			wantDomain: "domain.com",
			wantOk:     true,
		},
		{
			name:       "no separator",
			paystring:  "userdomain.com",
			wantUser:   "",
			wantDomain: "",
			wantOk:     false,
		},
		{
			name:       "multiple separators",
			paystring:  "user$middle$domain.com",
			wantUser:   "",
			wantDomain: "",
			wantOk:     false,
		},
		{
			name:       "empty string",
			paystring:  "",
			wantUser:   "",
			wantDomain: "",
			wantOk:     false,
		},
		{
			name:       "only separator",
			paystring:  "$",
			wantUser:   "",
			wantDomain: "",
			wantOk:     false, // Empty user and domain parts are invalid
		},
		{
			name:       "empty user",
			paystring:  "$domain.com",
			wantUser:   "",
			wantDomain: "",
			wantOk:     false,
		},
		{
			name:       "empty domain",
			paystring:  "user$",
			wantUser:   "",
			wantDomain: "",
			wantOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotDomain, gotOk := splitPaystring(tt.paystring)
			if gotUser != tt.wantUser {
				t.Errorf("splitPaystring() user = %v, want %v", gotUser, tt.wantUser)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("splitPaystring() domain = %v, want %v", gotDomain, tt.wantDomain)
			}
			if gotOk != tt.wantOk {
				t.Errorf("splitPaystring() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		paystring string
		wantErr   bool
	}{
		{
			name:      "valid paystring",
			paystring: "ifawoleesubiyi$paystring.crypto.com",
			wantErr:   false,
		},
		{
			name:      "valid simple paystring",
			paystring: "user$domain.com",
			wantErr:   false,
		},
		{
			name:      "empty paystring",
			paystring: "",
			wantErr:   true,
		},
		{
			name:      "missing dollar sign",
			paystring: "userdomain.com",
			wantErr:   true,
		},
		{
			name:      "multiple dollar signs",
			paystring: "user$extra$domain.com",
			wantErr:   true,
		},
		{
			name:      "empty user portion",
			paystring: "$domain.com",
			wantErr:   true,
		},
		{
			name:      "empty domain portion",
			paystring: "user$",
			wantErr:   true,
		},
		{
			name:      "domain without dot",
			paystring: "user$domain",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.paystring)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && p == nil {
				t.Error("New() returned nil payment for valid paystring")
			}
			if !tt.wantErr && p.GetPaystring() != tt.paystring {
				t.Errorf("New() paystring = %v, want %v", p.GetPaystring(), tt.paystring)
			}
		})
	}
}

func TestGetPaystring(t *testing.T) {
	paystring := "ifawoleesubiyi$paystring.crypto.com"
	p, err := New(paystring)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	got := p.GetPaystring()
	if got != paystring {
		t.Errorf("GetPaystring() = %v, want %v", got, paystring)
	}
}

func TestValidatePaystring(t *testing.T) {
	tests := []struct {
		name      string
		paystring string
		wantErr   bool
	}{
		{
			name:      "valid paystring",
			paystring: "ifawoleesubiyi$paystring.crypto.com",
			wantErr:   false,
		},
		{
			name:      "valid simple paystring",
			paystring: "alice$example.com",
			wantErr:   false,
		},
		{
			name:      "valid subdomain",
			paystring: "bob$pay.example.com",
			wantErr:   false,
		},
		{
			name:      "empty paystring",
			paystring: "",
			wantErr:   true,
		},
		{
			name:      "no dollar sign",
			paystring: "userexample.com",
			wantErr:   true,
		},
		{
			name:      "multiple dollar signs",
			paystring: "user$middle$example.com",
			wantErr:   true,
		},
		{
			name:      "empty user",
			paystring: "$example.com",
			wantErr:   true,
		},
		{
			name:      "empty domain",
			paystring: "user$",
			wantErr:   true,
		},
		{
			name:      "domain without dot",
			paystring: "user$localhost",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaystring(tt.paystring)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePaystring() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRedactedPaystring(t *testing.T) {
	tests := []struct {
		name      string
		paystring string
		want      string
	}{
		{
			name:      "long user portion",
			paystring: "ifawoleesubiyi$paystring.crypto.com",
			want:      "ifa***$paystring.crypto.com",
		},
		{
			name:      "medium user portion",
			paystring: "alice$example.com",
			want:      "ali***$example.com",
		},
		{
			name:      "short user portion (3 chars)",
			paystring: "bob$pay.example.com",
			want:      "bob***$pay.example.com",
		},
		{
			name:      "very short user portion (2 chars)",
			paystring: "ab$example.com",
			want:      "ab***$example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.paystring)
			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}

			got := p.GetRedactedPaystring()
			if got != tt.want {
				t.Errorf("GetRedactedPaystring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedactPaystring(t *testing.T) {
	tests := []struct {
		name      string
		paystring string
		want      string
	}{
		{
			name:      "long user portion",
			paystring: "ifawoleesubiyi$paystring.crypto.com",
			want:      "ifa***$paystring.crypto.com",
		},
		{
			name:      "medium user portion",
			paystring: "alice$example.com",
			want:      "ali***$example.com",
		},
		{
			name:      "short user portion (3 chars)",
			paystring: "bob$pay.example.com",
			want:      "bob***$pay.example.com",
		},
		{
			name:      "very short user portion (2 chars)",
			paystring: "ab$example.com",
			want:      "ab***$example.com",
		},
		{
			name:      "invalid format (no dollar sign)",
			paystring: "userexample.com",
			want:      "***",
		},
		{
			name:      "invalid format (multiple dollar signs)",
			paystring: "user$middle$example.com",
			want:      "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactPaystring(tt.paystring)
			if got != tt.want {
				t.Errorf("RedactPaystring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePayoutRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *PayoutRequest
		wantErr bool
	}{
		{
			name: "valid payout request",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: "test_payout",
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid recipient",
			req: &PayoutRequest{
				Recipient: "invalid",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    0,
				Currency:  "USD",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    -1050, // -$10.50 in cents
				Currency:  "USD",
			},
			wantErr: true,
		},
		{
			name: "empty currency",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "",
			},
			wantErr: true,
		},
		{
			name: "invalid currency code - lowercase",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "usd",
			},
			wantErr: true,
		},
		{
			name: "invalid currency code - too long",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USDD",
			},
			wantErr: true,
		},
		{
			name: "invalid currency code - too short",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "US",
			},
			wantErr: true,
		},
		{
			name: "amount exceeds maximum",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    200000000, // $2,000,000.00 in cents (exceeds max)
				Currency:  "USD",
			},
			wantErr: true,
		},
		{
			name: "valid payout request with maximum amount",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    100000000, // Exactly at the limit (MaxPayoutAmount)
				Currency:  "USD",
			},
			wantErr: false,
		},
		{
			name: "valid payout request with maximum amount EUR",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    100000000, // Exactly at the limit (MaxPayoutAmount) in cents
				Currency:  "EUR",
			},
			wantErr: false,
		},
		{
			name: "valid payout request with maximum amount GBP",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    100000000, // Exactly at the limit (MaxPayoutAmount) in pence
				Currency:  "GBP",
			},
			wantErr: false,
		},
		{
			name: "valid payout request with empty reference",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: "",
			},
			wantErr: false,
		},
		{
			name: "valid payout request with max length reference",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: strings.Repeat("a", 256), // exactly 256 characters
			},
			wantErr: false,
		},
		{
			name: "reference exceeds max length",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: strings.Repeat("a", 257), // 257 characters
			},
			wantErr: true,
		},
		{
			name: "reference contains control character - null",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: "test\x00payout",
			},
			wantErr: true,
		},
		{
			name: "reference contains control character - newline",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: "test\npayout",
			},
			wantErr: true,
		},
		{
			name: "reference contains control character - tab",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: "test\tpayout",
			},
			wantErr: true,
		},
		{
			name: "reference contains control character - DEL",
			req: &PayoutRequest{
				Recipient: "alice$example.com",
				Amount:    10050, // $100.50 in cents
				Currency:  "USD",
				Reference: "test\x7Fpayout",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePayoutRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePayoutRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutePayout(t *testing.T) {
	p, err := New("sender$example.com")
	if err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}

	tests := []struct {
		name    string
		req     *PayoutRequest
		wantErr bool
	}{
		{
			name: "successful payout",
			req: &PayoutRequest{
				Recipient: "recipient$example.com",
				Amount:    5000, // $50.00 in cents
				Currency:  "USD",
				Reference: "test",
			},
			wantErr: false,
		},
		{
			name: "invalid payout request",
			req: &PayoutRequest{
				Recipient: "invalid",
				Amount:    5000, // $50.00 in cents
				Currency:  "USD",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := p.ExecutePayout(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecutePayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result == nil {
					t.Error("ExecutePayout() returned nil result for valid request")
					return
				}
				if result.StatusCode != http.StatusOK {
					t.Errorf("ExecutePayout() StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
				}
				if result.Status != PayoutStatusCompleted {
					t.Errorf("ExecutePayout() Status = %s, want %s", result.Status, PayoutStatusCompleted)
				}
				if result.TransactionID == "" {
					t.Error("ExecutePayout() returned empty TransactionID")
				}
			}
		})
	}
}

func TestProcessPayout(t *testing.T) {
	p, err := New("sender$example.com")
	if err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}

	tests := []struct {
		name      string
		recipient string
		amount    int64
		currency  string
		wantErr   bool
	}{
		{
			name:      "successful payout",
			recipient: "recipient$example.com",
			amount:    10000, // $100.00 in cents
			currency:  "USD",
			wantErr:   false,
		},
		{
			name:      "invalid recipient - missing separator",
			recipient: "invalid",
			amount:    10000,
			currency:  "USD",
			wantErr:   true,
		},
		{
			name:      "invalid recipient - empty",
			recipient: "",
			amount:    10000,
			currency:  "USD",
			wantErr:   true,
		},
		{
			name:      "invalid amount - zero",
			recipient: "recipient$example.com",
			amount:    0,
			currency:  "USD",
			wantErr:   true,
		},
		{
			name:      "invalid amount - negative",
			recipient: "recipient$example.com",
			amount:    -100,
			currency:  "USD",
			wantErr:   true,
		},
		{
			name:      "invalid currency - empty",
			recipient: "recipient$example.com",
			amount:    10000,
			currency:  "",
			wantErr:   true,
		},
		{
			name:      "invalid currency - wrong format",
			recipient: "recipient$example.com",
			amount:    10000,
			currency:  "us",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := p.ProcessPayout(tt.recipient, tt.amount, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessPayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result == nil {
					t.Error("ProcessPayout() returned nil result for valid request")
					return
				}
				if result.StatusCode != http.StatusOK {
					t.Errorf("ProcessPayout() StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
				}
				if result.Status != PayoutStatusCompleted {
					t.Errorf("ProcessPayout() Status = %s, want %s", result.Status, PayoutStatusCompleted)
				}
				if result.TransactionID == "" {
					t.Error("ProcessPayout() returned empty TransactionID")
				}
			}
		})
	}
}
