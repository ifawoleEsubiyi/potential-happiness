package keystore

import (
	"testing"
)

func TestNew(t *testing.T) {
	k := New()
	if k == nil {
		t.Error("New() returned nil Keystore instance")
	}
}
