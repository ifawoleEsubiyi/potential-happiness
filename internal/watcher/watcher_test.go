package watcher

import (
	"testing"
)

func TestNew(t *testing.T) {
	w := New()
	if w == nil {
		t.Error("New() returned nil Watcher instance")
	}
}
