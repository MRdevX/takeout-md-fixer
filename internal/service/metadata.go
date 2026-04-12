package service

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"

	"takeout-md-fixer/internal/exif"
	"takeout-md-fixer/internal/pathkey"
	"takeout-md-fixer/internal/takeout"
)

// Checkpoint flush interval: balances crash safety vs disk I/O on large libraries.
const checkpointFlushEveryN = 25

// MetadataService exposes folder scan and metadata fix to the Wails frontend.
type MetadataService struct {
	fixMu     sync.Mutex
	fixCond   *sync.Cond
	fixActive bool
	fixPaused bool
	fixAbort  bool
}

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

// FixCheckpointAvailable reports whether a saved checkpoint exists for this folder (resume possible).
func (s *MetadataService) FixCheckpointAvailable(folderPath string) (bool, error) {
	return HasCheckpoint(folderPath)
}

// FixPause pauses the current fix between files.
func (s *MetadataService) FixPause() error {
	s.fixMu.Lock()
	defer s.fixMu.Unlock()
	if !s.fixActive {
		return ErrNoFixInProgress
	}
	s.fixPaused = true
	return nil
}

// FixResume continues a paused fix.
func (s *MetadataService) FixResume() error {
	s.fixMu.Lock()
	defer s.fixMu.Unlock()
	if !s.fixActive {
		return ErrNoFixInProgress
	}
	s.fixPaused = false
	if s.fixCond != nil {
		s.fixCond.Broadcast()
	}
	return nil
}

// FixAbort stops the current fix after the current file finishes (or immediately if waiting on pause).
func (s *MetadataService) FixAbort() error {
	s.fixMu.Lock()
	defer s.fixMu.Unlock()
	if !s.fixActive {
		return ErrNoFixInProgress
	}
	s.fixAbort = true
	s.fixPaused = false
	if s.fixCond != nil {
		s.fixCond.Broadcast()
	}
	return nil
}

func (s *MetadataService) beginFixSession() error {
	s.fixMu.Lock()
	defer s.fixMu.Unlock()
	if s.fixActive {
		return ErrFixAlreadyActive
	}
	if s.fixCond == nil {
		s.fixCond = sync.NewCond(&s.fixMu)
	}
	s.fixActive = true
	s.fixPaused = false
	s.fixAbort = false
	return nil
}

func (s *MetadataService) endFixSession() {
	s.fixMu.Lock()
	defer s.fixMu.Unlock()
	s.fixActive = false
	s.fixPaused = false
	s.fixAbort = false
	if s.fixCond != nil {
		s.fixCond.Broadcast()
	}
}

func (s *MetadataService) waitIfPaused() {
	s.fixMu.Lock()
	for s.fixPaused && !s.fixAbort {
		s.fixCond.Wait()
	}
	s.fixMu.Unlock()
}

func (s *MetadataService) shouldAbort() bool {
	s.fixMu.Lock()
	defer s.fixMu.Unlock()
	return s.fixAbort
}

// FixMetadata writes EXIF from sidecar JSON and optionally removes sidecar files.
// If resume is true, completed paths from .takeout-md-fixer-checkpoint.json are skipped.
//
// Work runs in a background goroutine so the Wails runtime can process other calls and events.
// The UI should listen for "fix-progress" and a final "fix-complete" event (see FixComplete).
// A non-nil error means the job could not be started (e.g. another fix is active).
func (s *MetadataService) FixMetadata(folderPath string, deleteJsonSidecars bool, resume bool) error {
	if err := s.beginFixSession(); err != nil {
		return err
	}
	go s.runFix(folderPath, deleteJsonSidecars, resume)
	return nil
}

func (s *MetadataService) runFix(folderPath string, deleteJsonSidecars bool, resume bool) {
	defer s.endFixSession()

	app := application.Get()
	emitComplete := func(c FixComplete) {
		app.Event.Emit("fix-complete", c)
	}

	scanResult, err := takeout.ScanFolder(folderPath)
	if err != nil {
		emitComplete(FixComplete{Error: err.Error()})
		return
	}

	if !resume {
		if err := clearCheckpoint(folderPath); err != nil {
			emitComplete(FixComplete{Error: err.Error()})
			return
		}
	}

	completedByNorm, err := loadCheckpointState(folderPath, deleteJsonSidecars)
	if err != nil {
		emitComplete(FixComplete{Error: err.Error()})
		return
	}
	if completedByNorm == nil {
		completedByNorm = make(map[string]string)
	}

	writer, err := exif.NewWriter()
	if err != nil {
		emitComplete(FixComplete{Error: err.Error()})
		return
	}
	defer writer.Close()

	result := &takeout.FixResult{
		Total:   len(scanResult.Files),
		Resumed: resume && len(completedByNorm) > 0,
	}

	sinceFlush := 0
	flushPending := func() error {
		if sinceFlush == 0 {
			return nil
		}
		sinceFlush = 0
		return writeCheckpoint(folderPath, deleteJsonSidecars, completedByNorm)
	}

	for i, mf := range scanResult.Files {
		s.waitIfPaused()
		if s.shouldAbort() {
			result.Aborted = true
			break
		}

		progress := takeout.FixProgress{
			Current: i + 1,
			Total:   result.Total,
			File:    mf.Name,
		}

		if CheckpointContains(completedByNorm, mf.Path) {
			progress.Status = "resumed"
			result.Success++
			app.Event.Emit("fix-progress", progress)
			continue
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

		// Each media path is updated independently (including Live Photo motion files that borrow JSON from the still).
		if err := writer.WriteMetadata(mf.Path, meta); err != nil {
			progress.Status = "error"
			result.Failed++
			app.Event.Emit("fix-progress", progress)
			continue
		}

		if deleteJsonSidecars && mf.JsonPath != "" {
			for _, p := range takeout.SidecarDeletionPaths(mf.Path, mf.JsonPath) {
				if err := os.Remove(p); err != nil {
					if os.IsNotExist(err) {
						continue
					}
					result.JsonDeleteFailed++
				} else {
					result.JsonDeleted++
				}
			}
		}

		abs, err := filepath.Abs(mf.Path)
		if err == nil {
			completedByNorm[pathkey.Normalize(abs)] = abs
			sinceFlush++
			if sinceFlush >= checkpointFlushEveryN {
				if err := writeCheckpoint(folderPath, deleteJsonSidecars, completedByNorm); err != nil {
					emitComplete(FixComplete{Result: result, Error: err.Error()})
					return
				}
				sinceFlush = 0
			}
		}

		progress.Status = "success"
		result.Success++
		app.Event.Emit("fix-progress", progress)

		if s.shouldAbort() {
			result.Aborted = true
			break
		}
	}

	if err := flushPending(); err != nil {
		emitComplete(FixComplete{Result: result, Error: err.Error()})
		return
	}

	if result.Aborted {
		emitComplete(FixComplete{Result: result})
		return
	}

	if err := clearCheckpoint(folderPath); err != nil {
		emitComplete(FixComplete{Result: result, Error: err.Error()})
		return
	}

	emitComplete(FixComplete{Result: result})
}
