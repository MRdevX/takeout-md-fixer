package service

import "takeout-md-fixer/internal/takeout"

// FixComplete is emitted on the "fix-complete" application event when FixMetadata finishes.
type FixComplete struct {
	Result *takeout.FixResult `json:"result,omitempty"`
	Error  string             `json:"error,omitempty"`
}
