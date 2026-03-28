package takeout

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScanFolder walks root recursively and lists media files with optional JSON sidecars.
func ScanFolder(root string) (*ScanResult, error) {
	if root == "" {
		return nil, fmt.Errorf("no folder path provided")
	}

	result := &ScanResult{FolderPath: root}

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

	return result, nil
}
