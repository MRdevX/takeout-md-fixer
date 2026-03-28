package takeout

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFolder_orphanJsonAndSkipsMetadata(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "a.jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	sidecar := filepath.Join(dir, "a.jpg.json")
	if err := os.WriteFile(sidecar, []byte(`{"photoTakenTime":{"timestamp":"0"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	orphan := filepath.Join(dir, "lonely.json")
	if err := os.WriteFile(orphan, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	meta := filepath.Join(dir, "metadata.json")
	if err := os.WriteFile(meta, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := ScanFolder(dir)
	if err != nil {
		t.Fatal(err)
	}
	if res.TotalMedia != 1 {
		t.Fatalf("TotalMedia got %d want 1", res.TotalMedia)
	}
	if res.WithJson != 1 || res.WithoutJson != 0 {
		t.Fatalf("with/without json %d %d", res.WithJson, res.WithoutJson)
	}
	if res.OrphanJson != 1 {
		t.Fatalf("OrphanJson got %d want 1 (lonely.json only; metadata.json skipped)", res.OrphanJson)
	}
}
