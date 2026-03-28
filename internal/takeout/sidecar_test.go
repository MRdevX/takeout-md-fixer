package takeout

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSidecarPath_jpgJsonAndPhotoJson(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	jpgJSON := filepath.Join(dir, "photo.jpg.json")
	if err := os.WriteFile(jpgJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	// photo.jpg.json wins over photo.json when both exist
	onlyJSON := filepath.Join(dir, "photo.json")
	if err := os.WriteFile(onlyJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != jpgJSON {
		t.Fatalf("got %q want %q", got, jpgJSON)
	}
}

func TestSidecarPath_photoJsonOnly(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	onlyJSON := filepath.Join(dir, "photo.json")
	if err := os.WriteFile(onlyJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != onlyJSON {
		t.Fatalf("got %q want %q", got, onlyJSON)
	}
}

func TestSidecarPath_supplementalMetadata(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	sup := filepath.Join(dir, "photo.jpg.supplemental-metadata.json")
	if err := os.WriteFile(sup, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != sup {
		t.Fatalf("got %q want %q", got, sup)
	}
}

func TestSidecarPath_supplementalTruncatedInDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "photo.jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	trunc := filepath.Join(dir, "photo.jpg.supplemental-metad.json")
	if err := os.WriteFile(trunc, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != trunc {
		t.Fatalf("got %q want %q", got, trunc)
	}
}

func TestSidecarPath_duplicateOriginalJsonFirst(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "IMG(1).jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	origJSON := filepath.Join(dir, "IMG.json")
	if err := os.WriteFile(origJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	altJSON := filepath.Join(dir, "IMG.jpg(1).json")
	if err := os.WriteFile(altJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != origJSON {
		t.Fatalf("got %q want %q", got, origJSON)
	}
}

func TestSidecarPath_duplicateJpgNumberedJson(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "IMG(1).jpg")
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	altJSON := filepath.Join(dir, "IMG.jpg(1).json")
	if err := os.WriteFile(altJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != altJSON {
		t.Fatalf("got %q want %q", got, altJSON)
	}
}

func TestSidecarPath_longBasenameTruncatedJson(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	// 43 'a' + .jpg => basename runes 47 (≥47 triggers long-name handling)
	base := strings.Repeat("a", 43) + ".jpg"
	media := filepath.Join(dir, base)
	if err := os.WriteFile(media, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	truncName := prefixRunes(base, takeoutBasenameTruncateRunes) + ".json"
	truncPath := filepath.Join(dir, truncName)
	if err := os.WriteFile(truncPath, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(media); got != truncPath {
		t.Fatalf("got %q want %q", got, truncPath)
	}
}

func TestPrefixRunes(t *testing.T) {
	t.Parallel()
	if got, want := prefixRunes("abcdef", 3), "abc"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if got, want := prefixRunes("éclair", 2), "éc"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
