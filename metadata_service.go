package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var mediaExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".heic": true, ".webp": true, ".bmp": true, ".tiff": true,
	".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	".3gp": true, ".m4v": true,
}

type TakeoutMeta struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	PhotoTakenTime TimeStamp `json:"photoTakenTime"`
	GeoData        GeoData   `json:"geoData"`
	GeoDataExif    GeoData   `json:"geoDataExif"`
}

type TimeStamp struct {
	Timestamp string `json:"timestamp"`
	Formatted string `json:"formatted"`
}

type GeoData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type MediaFile struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	JsonPath string `json:"jsonPath"`
	HasJson  bool   `json:"hasJson"`
	Status   string `json:"status"`
	Error    string `json:"error"`
}

type ScanResult struct {
	FolderPath  string      `json:"folderPath"`
	TotalMedia  int         `json:"totalMedia"`
	WithJson    int         `json:"withJson"`
	WithoutJson int         `json:"withoutJson"`
	Files       []MediaFile `json:"files"`
}

type FixProgress struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	File    string `json:"file"`
	Status  string `json:"status"`
}

type FixResult struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
}

type MetadataService struct{}

func (s *MetadataService) SelectFolder() (string, error) {
	app := application.Get()
	dialog := app.Dialog.OpenFile()
	dialog.CanChooseDirectories(true)
	dialog.CanChooseFiles(false)
	dialog.SetTitle("Select Google Takeout Folder")
	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s *MetadataService) ScanFolder(folderPath string) (*ScanResult, error) {
	if folderPath == "" {
		return nil, fmt.Errorf("no folder path provided")
	}

	result := &ScanResult{FolderPath: folderPath}

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
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
		if !mediaExtensions[ext] {
			return nil
		}

		mf := MediaFile{
			Path:   path,
			Name:   info.Name(),
			Status: "pending",
		}

		jsonPath := findJsonSidecar(path)
		if jsonPath != "" {
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

func findJsonSidecar(mediaPath string) string {
	candidates := []string{
		mediaPath + ".json",
		mediaPath + ".supplemental-metadata.json",
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	// Handle Google's duplicate naming: MOVIE(1).mp4 -> MOVIE.mp4(1).json
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

func (s *MetadataService) FixMetadata(folderPath string) (*FixResult, error) {
	scanResult, err := s.ScanFolder(folderPath)
	if err != nil {
		return nil, err
	}

	app := application.Get()
	writer, err := NewExifWriter()
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	result := &FixResult{Total: len(scanResult.Files)}

	for i, mf := range scanResult.Files {
		progress := FixProgress{
			Current: i + 1,
			Total:   result.Total,
			File:    mf.Name,
		}

		if !mf.HasJson {
			progress.Status = "skipped"
			result.Skipped++
			app.Event.Emit("fix-progress", progress)
			continue
		}

		meta, err := parseTakeoutJson(mf.JsonPath)
		if err != nil {
			progress.Status = "error"
			result.Failed++
			app.Event.Emit("fix-progress", progress)
			continue
		}

		if err := writer.WriteMetadata(mf.Path, meta); err != nil {
			progress.Status = "error"
			result.Failed++
			app.Event.Emit("fix-progress", progress)
			continue
		}

		progress.Status = "success"
		result.Success++
		app.Event.Emit("fix-progress", progress)
	}

	return result, nil
}

func parseTakeoutJson(jsonPath string) (*TakeoutMeta, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", jsonPath, err)
	}

	var meta TakeoutMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", jsonPath, err)
	}

	// Prefer geoDataExif over geoData when available and non-zero
	if meta.GeoDataExif.Latitude != 0 || meta.GeoDataExif.Longitude != 0 {
		meta.GeoData = meta.GeoDataExif
	}

	return &meta, nil
}
