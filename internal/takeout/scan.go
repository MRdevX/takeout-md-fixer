package takeout

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// ScanFolder walks root recursively and lists media files with optional JSON sidecars,
// then counts JSON files that were not linked to any media (orphans), excluding known non-sidecar names.
// Filesystem errors on individual paths (permissions, broken symlinks, etc.) are skipped so the rest of the tree still scans.
func ScanFolder(root string) (*ScanResult, error) {
	if root == "" {
		return nil, fmt.Errorf("no folder path provided")
	}

	result := &ScanResult{FolderPath: root}
	linkedJSON := make(map[string]struct{})
	var jsonCandidates []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, errInfo := d.Info()
		if errInfo != nil {
			return nil
		}

		// macOS AppleDouble resource forks (._filename) masquerade as media + .json; JSON is binary, not Takeout.
		if strings.HasPrefix(info.Name(), "._") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if MediaExtensions[ext] {
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
		}

		if strings.HasSuffix(strings.ToLower(path), ".json") {
			jsonCandidates = append(jsonCandidates, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning folder: %w", err)
	}

	for _, path := range jsonCandidates {
		b := strings.ToLower(filepath.Base(path))
		if b == "metadata.json" || b == "print-subscriptions.json" {
			continue
		}
		k := normalizePathKey(path)
		if k == "" {
			continue
		}
		if _, ok := linkedJSON[k]; ok {
			continue
		}
		result.OrphanJson++
	}

	return result, nil
}

// ListAlbumMetadataJSON walks root and returns paths to Google Takeout per-album metadata.json files
// (album title only; not photo sidecars). Skips AppleDouble "._*" files.
func ListAlbumMetadataJSON(root string) ([]string, error) {
	if root == "" {
		return nil, nil
	}
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, "._") {
			return nil
		}
		if strings.EqualFold(name, "metadata.json") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func normalizePathKey(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	return strings.ToLower(filepath.Clean(abs))
}
