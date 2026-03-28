package exif

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAugmentedPath_preservesPATHFirstAndDedupes(t *testing.T) {
	t.Parallel()
	sep := string(os.PathListSeparator)
	first := filepath.Join(t.TempDir(), "first")
	second := filepath.Join(t.TempDir(), "second")
	base := first + sep + second
	got := augmentedPath(base)
	if !strings.HasPrefix(got, base) {
		t.Fatalf("expected PATH to keep original entries first: %q", got)
	}
	// duplicate dir in extras should not repeat first
	if strings.Count(got, first) > 1 {
		t.Fatalf("duplicate path segment: %q", got)
	}
}

func TestAugmentedPath_includesExpectedExtras(t *testing.T) {
	t.Parallel()
	got := augmentedPath("/usr/bin")
	switch runtime.GOOS {
	case "darwin":
		for _, want := range []string{"/opt/homebrew/bin", "/usr/local/bin", "/opt/local/bin"} {
			if !strings.Contains(got, want) {
				t.Errorf("expected augmented PATH to contain %q, got %q", want, got)
			}
		}
	case "linux":
		if !strings.Contains(got, "/usr/local/bin") {
			t.Errorf("expected augmented PATH to contain /usr/local/bin, got %q", got)
		}
	}
}

func writeFakeExiftool(t *testing.T, dir string) string {
	t.Helper()
	name := "exiftool"
	if runtime.GOOS == "windows" {
		name = "exiftool.exe"
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		// Windows does not use shebang execute bits the same way; empty exe still discoverable if .exe exists.
		if err := os.Chmod(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func TestResolveExiftoolPath_findsOnPATH(t *testing.T) {
	tmp := t.TempDir()
	want := writeFakeExiftool(t, tmp)
	t.Setenv("PATH", tmp)
	got, err := ResolveExiftoolPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(want) {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestResolveExiftoolPath_findsViaAugmentedDirsWhenNotOnPATH(t *testing.T) {
	tmp := t.TempDir()
	want := writeFakeExiftool(t, tmp)
	t.Setenv("PATH", filepath.Join(tmp, "missing-bin")) // no exiftool here
	extraExiftoolPathDirsOverride = func() []string { return []string{tmp} }
	defer func() { extraExiftoolPathDirsOverride = nil }()
	got, err := ResolveExiftoolPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(want) {
		t.Fatalf("got %q want %q", got, want)
	}
}
