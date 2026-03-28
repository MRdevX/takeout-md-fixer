package service

import (
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"

	"takeout-md-fixer/internal/exif"
	"takeout-md-fixer/internal/takeout"
)

// MetadataService exposes folder scan and metadata fix to the Wails frontend.
type MetadataService struct{}

// ExiftoolStatus is returned by ExiftoolCheck for the UI.
type ExiftoolStatus struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Path    string `json:"path,omitempty"`
}

// ExiftoolCheck reports whether ExifTool is available (PATH plus common install locations).
func (s *MetadataService) ExiftoolCheck() ExiftoolStatus {
	path, err := exif.ResolveExiftoolPath()
	if err != nil {
		return ExiftoolStatus{OK: false, Message: err.Error()}
	}
	return ExiftoolStatus{OK: true, Path: path}
}

// SelectFolder opens a native directory picker.
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

// ScanFolder lists media files under folderPath and their JSON sidecars.
func (s *MetadataService) ScanFolder(folderPath string) (*takeout.ScanResult, error) {
	return takeout.ScanFolder(folderPath)
}

// FixMetadata writes EXIF from sidecar JSON and optionally removes sidecar files.
func (s *MetadataService) FixMetadata(folderPath string, deleteJsonSidecars bool) (*takeout.FixResult, error) {
	scanResult, err := takeout.ScanFolder(folderPath)
	if err != nil {
		return nil, err
	}

	app := application.Get()
	writer, err := exif.NewWriter()
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	result := &takeout.FixResult{Total: len(scanResult.Files)}

	for i, mf := range scanResult.Files {
		progress := takeout.FixProgress{
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

		meta, err := takeout.ParseMetadataFile(mf.JsonPath)
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

		if deleteJsonSidecars && mf.JsonPath != "" {
			if err := os.Remove(mf.JsonPath); err != nil {
				result.JsonDeleteFailed++
			} else {
				result.JsonDeleted++
			}
		}

		progress.Status = "success"
		result.Success++
		app.Event.Emit("fix-progress", progress)
	}

	return result, nil
}
