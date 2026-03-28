package exif

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var resolveMu sync.Mutex

// ResolveExiftoolPath finds the ExifTool binary on PATH augmented with common install locations
// (GUI apps on macOS often have a minimal PATH).
func ResolveExiftoolPath() (string, error) {
	name := exiftoolBinaryName()
	resolveMu.Lock()
	defer resolveMu.Unlock()
	oldPath := os.Getenv("PATH")
	defer func() { _ = os.Setenv("PATH", oldPath) }()
	if err := os.Setenv("PATH", augmentedPath(oldPath)); err != nil {
		return "", fmt.Errorf("could not search for ExifTool: %w", err)
	}
	p, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("ExifTool was not found. Install it from https://exiftool.org/ (e.g. Homebrew or the official package) " +
			"so `exiftool` is in a standard location or on your system PATH")
	}
	return p, nil
}

func exiftoolBinaryName() string {
	if runtime.GOOS == "windows" {
		return "exiftool.exe"
	}
	return "exiftool"
}

// augmentedPath returns PATH with OS-specific directories appended (deduplicated) so exec.LookPath
// finds ExifTool when installed outside the GUI process PATH (e.g. /opt/homebrew/bin on macOS).
func augmentedPath(existing string) string {
	sep := string(os.PathListSeparator)
	parts := filepath.SplitList(existing)
	extras := extraExiftoolPathDirs()
	seen := make(map[string]bool, len(parts)+len(extras))
	var out []string
	for _, p := range parts {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	for _, p := range extras {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return strings.Join(out, sep)
}

// extraExiftoolPathDirsOverride is set by tests to simulate installs outside PATH.
var extraExiftoolPathDirsOverride func() []string

func extraExiftoolPathDirs() []string {
	if extraExiftoolPathDirsOverride != nil {
		return extraExiftoolPathDirsOverride()
	}
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/opt/homebrew/bin",
			"/usr/local/bin",
			"/opt/local/bin",
		}
	case "windows":
		return windowsExtraPathDirs()
	case "linux", "freebsd", "openbsd", "netbsd":
		return []string{"/usr/local/bin", "/usr/bin"}
	default:
		return []string{"/usr/local/bin"}
	}
}

func windowsExtraPathDirs() []string {
	var out []string
	addIfDir := func(dir string) {
		if dir == "" {
			return
		}
		if st, err := os.Stat(dir); err == nil && st.IsDir() {
			out = append(out, dir)
		}
	}
	if pf := os.Getenv("ProgramFiles"); pf != "" {
		addIfDir(filepath.Join(pf, "exiftool"))
	}
	if pfx86 := os.Getenv("ProgramFiles(x86)"); pfx86 != "" {
		addIfDir(filepath.Join(pfx86, "exiftool"))
	}
	addIfDir(`C:\ProgramData\chocolatey\bin`)
	return out
}
