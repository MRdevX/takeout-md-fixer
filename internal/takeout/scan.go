package takeout

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScanFolder walks root recursively and lists media files with optional JSON sidecars,
// then counts JSON files that were not linked to any media (orphans), excluding known non-sidecar names.
func ScanFolder(root string) (*ScanResult, error) {
	if root == "" {
		return nil, fmt.Errorf("no folder path provided")
	}

	result := &ScanResult{FolderPath: root}
	linkedJSON := make(map[string]struct{})

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// macOS AppleDouble resource forks (._filename) masquerade as media + .json; JSON is binary, not Takeout.
		if strings.HasPrefix(info.Name(), "._") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !MediaExtensions[ext] {
			return nil
		}

		mf := MediaFile{
			Path:   path,
			Name:   info.Name(),
			Status: "pending",
		}

		if jsonPath := SidecarPath(path); jsonPath != "" {
			mf.JsonPath = jsonPath
			mf.HasJson = true
			result.WithJson++
			if k := normalizePathKey(jsonPath); k != "" {
				linkedJSON[k] = struct{}{}
			}
		} else {
			result.WithoutJson++
		}

		result.Files = append(result.Files, mf)
		result.TotalMedia++
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning folder: %w", err)
	}

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), "._") {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".json") {
			return nil
		}
		b := strings.ToLower(filepath.Base(path))
		if b == "metadata.json" || b == "print-subscriptions.json" {
			return nil
		}
		k := normalizePathKey(path)
		if k == "" {
			return nil
		}
		if _, ok := linkedJSON[k]; ok {
			return nil
		}
		result.OrphanJson++
		return nil
	})

	return result, nil
}

func normalizePathKey(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	return strings.ToLower(filepath.Clean(abs))
}
