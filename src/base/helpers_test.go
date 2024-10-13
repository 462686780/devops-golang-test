package base

import (
	"testing"
)

func TestNewRequestID(t *testing.T) {
	in, s := NewRequestID()
	if in == 0 {
		t.Error(`NewRequestID failed`)
	}
	if s == "" {
		t.Error(`NewRequestID failed`)
	}
}
