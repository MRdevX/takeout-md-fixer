package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"takeout-md-fixer/internal/pathkey"
)

const checkpointFileName = ".takeout-md-fixer-checkpoint.json"

// checkpointData is written next to the Takeout folder being processed.
type checkpointData struct {
	Version            int      `json:"version"`
	FolderPath         string   `json:"folderPath"`
	DeleteJsonSidecars bool     `json:"deleteJsonSidecars"`
	CompletedPaths     []string `json:"completedPaths"`
}

func checkpointPath(folder string) string {
	return filepath.Join(folder, checkpointFileName)
}

// HasCheckpoint reports whether a checkpoint file exists for the folder.
func HasCheckpoint(folder string) (bool, error) {
	p := checkpointPath(folder)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// loadCheckpointState returns normalized-key -> absolute path for completed files, or empty map if none.
func loadCheckpointState(folder string, deleteJsonSidecars bool) (map[string]string, error) {
	p := checkpointPath(folder)
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var c checkpointData
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("checkpoint file: %w", err)
	}
	if c.Version != 1 {
		return nil, fmt.Errorf("unsupported checkpoint version %d", c.Version)
	}
	folderAbs, err := filepath.Abs(folder)
	if err != nil {
		return nil, err
	}
	cpFolder, err := filepath.Abs(c.FolderPath)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(filepath.Clean(cpFolder), filepath.Clean(folderAbs)) {
		return nil, errors.New("checkpoint folder path does not match")
	}
	if c.DeleteJsonSidecars != deleteJsonSidecars {
		return nil, errors.New("checkpoint was created with a different sidecar deletion setting; run without resume or delete the checkpoint file")
	}
	out := make(map[string]string, len(c.CompletedPaths))
	for _, path := range c.CompletedPaths {
		if path == "" {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		k := pathkey.Normalize(abs)
		out[k] = abs
	}
	return out, nil
}

// writeCheckpoint saves byNorm keys (normalized lookup) mapped to absolute paths for JSON.
func writeCheckpoint(folder string, deleteJson bool, byNorm map[string]string) error {
	folderAbs, err := filepath.Abs(folder)
	if err != nil {
		return err
	}
	paths := make([]string, 0, len(byNorm))
	for _, abs := range byNorm {
		paths = append(paths, abs)
	}
	sort.Strings(paths)
	c := checkpointData{
		Version:            1,
		FolderPath:         folderAbs,
		DeleteJsonSidecars: deleteJson,
		CompletedPaths:     paths,
	}
	raw, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	dest := checkpointPath(folder)
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, dest)
}

func clearCheckpoint(folder string) error {
	p := checkpointPath(folder)
	err := os.Remove(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// CheckpointContains reports whether mediaPath is marked completed.
func CheckpointContains(completed map[string]string, mediaPath string) bool {
	if completed == nil {
		return false
	}
	_, ok := completed[pathkey.Normalize(mediaPath)]
	return ok
}
