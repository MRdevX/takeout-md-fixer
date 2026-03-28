package exif

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const EnvExiftoolPath = "TAKEOUT_EXIFTOOL_PATH"

// ResolveExiftoolPath finds the ExifTool binary using TAKEOUT_EXIFTOOL_PATH if set, otherwise PATH.
func ResolveExiftoolPath() (string, error) {
	if p := os.Getenv(EnvExiftoolPath); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return filepath.Abs(p)
		}
		return "", fmt.Errorf("%s does not point to a valid file: %q", EnvExiftoolPath, p)
	}

	name := "exiftool"
	if runtime.GOOS == "windows" {
		name = "exiftool.exe"
	}
	p, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("ExifTool was not found. Install it from https://exiftool.org/ and ensure it is on your PATH (or set %s to the full path of the binary)", EnvExiftoolPath)
	}
	return p, nil
}
