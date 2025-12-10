package attest

import (
	"testing"
)

func TestNew(t *testing.T) {
	a := New()
	if a == nil {
		t.Error("New() returned nil Attest instance")
	}
}
