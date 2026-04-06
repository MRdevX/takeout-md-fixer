package service

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	t.Parallel()
	wrappedNo := fmt.Errorf("ctx: %w", ErrNoFixInProgress)
	if !errors.Is(wrappedNo, ErrNoFixInProgress) {
		t.Fatal("errors.Is should unwrap ErrNoFixInProgress")
	}
	wrappedBusy := fmt.Errorf("ctx: %w", ErrFixAlreadyActive)
	if !errors.Is(wrappedBusy, ErrFixAlreadyActive) {
		t.Fatal("errors.Is should unwrap ErrFixAlreadyActive")
	}
}
