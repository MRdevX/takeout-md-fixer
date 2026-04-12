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

func TestScanFolder_livePhotoJpgAndMp4BothHaveJson(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	jpg := filepath.Join(dir, "IMG_5175.JPG")
	mp4 := filepath.Join(dir, "IMG_5175.MP4")
	sup := filepath.Join(dir, "IMG_5175.JPG.supplemental-metadata.json")
	meta := `{"photoTakenTime":{"timestamp":"1701460188"},"geoData":{"latitude":1,"longitude":2}}`
	for _, p := range []string{jpg, mp4} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(sup, []byte(meta), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := ScanFolder(dir)
	if err != nil {
		t.Fatal(err)
	}
	if res.TotalMedia != 2 {
		t.Fatalf("TotalMedia got %d want 2", res.TotalMedia)
	}
	if res.WithJson != 2 || res.WithoutJson != 0 {
		t.Fatalf("with/without json got %d %d want 2 0", res.WithJson, res.WithoutJson)
	}
	byName := make(map[string]bool)
	for _, f := range res.Files {
		byName[f.Name] = f.HasJson
	}
	if !byName["IMG_5175.JPG"] || !byName["IMG_5175.MP4"] {
		t.Fatalf("expected both JPG and MP4 HasJson true: %#v", byName)
	}
}

func TestListAlbumMetadataJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	sub := filepath.Join(dir, "Google Photos", "My Album")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	meta := filepath.Join(sub, "metadata.json")
	if err := os.WriteFile(meta, []byte(`{"title":"My Album"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	photoMeta := filepath.Join(sub, "IMG_1.JPG.json")
	if err := os.WriteFile(photoMeta, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	paths, err := ListAlbumMetadataJSON(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 || filepath.Clean(paths[0]) != filepath.Clean(meta) {
		t.Fatalf("got %#v want single album metadata.json", paths)
	}
}
