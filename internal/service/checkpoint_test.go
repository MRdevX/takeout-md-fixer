package service

import (
	"path/filepath"
	"testing"

	"takeout-md-fixer/internal/pathkey"
)

func TestCheckpointRoundTrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	m1 := filepath.Join(dir, "a.jpg")
	m2 := filepath.Join(dir, "b.jpg")
	byNorm := map[string]string{
		pathkey.Normalize(m1): m1,
		pathkey.Normalize(m2): m2,
	}
	if err := writeCheckpoint(dir, false, byNorm); err != nil {
		t.Fatal(err)
	}
	loaded, err := loadCheckpointState(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("want 2 entries, got %d", len(loaded))
	}
	if !CheckpointContains(loaded, m1) || !CheckpointContains(loaded, m2) {
		t.Fatal("missing paths")
	}
	if err := clearCheckpoint(dir); err != nil {
		t.Fatal(err)
	}
	ok, err := HasCheckpoint(dir)
	if err != nil || ok {
		t.Fatalf("checkpoint should be gone: ok=%v err=%v", ok, err)
	}
}

func TestCheckpointDeleteJsonMismatch(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	byNorm := map[string]string{pathkey.Normalize(filepath.Join(dir, "a.jpg")): filepath.Join(dir, "a.jpg")}
	if err := writeCheckpoint(dir, true, byNorm); err != nil {
		t.Fatal(err)
	}
	_, err := loadCheckpointState(dir, false)
	if err == nil {
		t.Fatal("expected error when deleteJsonSidecars differs")
	}
}
