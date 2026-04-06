package pathkey

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "Photo.JPG")
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.ToLower(filepath.Clean(abs))
	if got := Normalize(p); got != want {
		t.Fatalf("Normalize(%q) = %q, want %q", p, got, want)
	}
}
