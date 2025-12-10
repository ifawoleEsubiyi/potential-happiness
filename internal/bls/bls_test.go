package bls

import (
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	if b == nil {
		t.Error("New() returned nil BLS instance")
	}
}
