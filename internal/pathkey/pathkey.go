// Package pathkey provides a single canonical string key for comparing filesystem paths
// across packages (scanning, checkpoints). Keys are lowercased cleaned paths, with
// filepath.Abs applied when possible.
package pathkey

import (
	"path/filepath"
	"strings"
)

// Normalize returns a stable comparison key for p. If Abs fails, it falls back to
// strings.ToLower(filepath.Clean(p)).
func Normalize(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return strings.ToLower(filepath.Clean(p))
	}
	return strings.ToLower(filepath.Clean(abs))
}
