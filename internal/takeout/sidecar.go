package takeout

import (
	"os"
	"path/filepath"
	"strings"
)

// SidecarPath returns the path to the JSON metadata file for a media file, or "" if none.
func SidecarPath(mediaPath string) string {
	candidates := []string{
		mediaPath + ".json",
		mediaPath + ".supplemental-metadata.json",
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	// Duplicate naming: MOVIE(1).mp4 -> MOVIE.mp4(1).json
	dir := filepath.Dir(mediaPath)
	base := filepath.Base(mediaPath)
	ext := filepath.Ext(base)
	nameNoExt := strings.TrimSuffix(base, ext)

	if idx := strings.LastIndex(nameNoExt, "("); idx > 0 {
		suffix := nameNoExt[idx:]
		originalName := nameNoExt[:idx]
		candidate := filepath.Join(dir, originalName+ext+suffix+".json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}
