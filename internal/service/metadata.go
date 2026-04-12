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

// FixCheckpointAvailable reports whether a saved checkpoint exists for this folder (UI uses this for auto-resume and a short notice).
func (s *MetadataService) FixCheckpointAvailable(folderPath string) (bool, error) {
	return hasCheckpoint(folderPath)
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
// Sidecar deletion runs once at the end: each JSON is removed only after every media file
// that referenced it has been written (including items skipped via resume checkpoint).
// If resume is true, completed paths from .takeout-md-fixer-checkpoint.json are skipped (the UI sets this when that file exists).
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
	emitProgress := func(p takeout.FixProgress) {
		app.Event.Emit("fix-progress", p)
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
	defer func() { _ = writer.Close() }()

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

		if checkpointContains(completedByNorm, mf.Path) {
			progress.Status = "resumed"
			result.Success++
			emitProgress(progress)
			continue
		}

		if !mf.HasJson {
			progress.Status = "skipped"
			result.Skipped++
			emitProgress(progress)
			continue
		}

		meta, err := takeout.ParseMetadataFile(mf.JsonPath)
		if err != nil {
			progress.Status = "error"
			result.Failed++
			emitProgress(progress)
			continue
		}

		// Each media path is updated independently (including Live Photo motion files that borrow JSON from the still).
		if err := writer.WriteMetadata(mf.Path, meta); err != nil {
			progress.Status = "error"
			result.Failed++
			emitProgress(progress)
			continue
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
		emitProgress(progress)

		if s.shouldAbort() {
			result.Aborted = true
			break
		}
	}

	if err := flushPending(); err != nil {
		emitComplete(FixComplete{Result: result, Error: err.Error()})
		return
	}

	if deleteJsonSidecars && !result.Aborted {
		deleteSharedSidecarsAfterFix(scanResult, completedByNorm, result)
		deleteAlbumMetadataJSONFiles(folderPath, result)
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

// deleteSharedSidecarsAfterFix removes Takeout JSON for each sidecar only after every media file
// that used that sidecar has been written successfully (including paths already in the checkpoint from resume).
func deleteSharedSidecarsAfterFix(scanResult *takeout.ScanResult, completedByNorm map[string]string, result *takeout.FixResult) {
	byJSON := make(map[string][]takeout.MediaFile)
	for _, mf := range scanResult.Files {
		if !mf.HasJson || mf.JsonPath == "" {
			continue
		}
		k := pathkey.Normalize(mf.JsonPath)
		byJSON[k] = append(byJSON[k], mf)
	}

	removed := make(map[string]struct{})
	for _, group := range byJSON {
		allDone := true
		for _, mf := range group {
			if !checkpointContains(completedByNorm, mf.Path) {
				allDone = false
				break
			}
		}
		if !allDone {
			continue
		}

		jsonPath := group[0].JsonPath
		mediaPaths := make([]string, len(group))
		for i := range group {
			mediaPaths[i] = group[i].Path
		}
		rep := takeout.SidecarCleanupRepresentative(mediaPaths, jsonPath)
		if rep == "" {
			continue
		}
		for _, p := range takeout.SidecarCleanupPaths(rep, jsonPath) {
			cp := filepath.Clean(p)
			if _, ok := removed[cp]; ok {
				continue
			}
			if err := os.Remove(p); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				result.JsonDeleteFailed++
				continue
			}
			removed[cp] = struct{}{}
			result.JsonDeleted++
		}
	}
}

// deleteAlbumMetadataJSONFiles removes per-album metadata.json files (album title) under the Takeout tree.
func deleteAlbumMetadataJSONFiles(folderPath string, result *takeout.FixResult) {
	paths, err := takeout.ListAlbumMetadataJSON(folderPath)
	if err != nil {
		return
	}
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			result.JsonDeleteFailed++
			continue
		}
		result.JsonDeleted++
	}
}
