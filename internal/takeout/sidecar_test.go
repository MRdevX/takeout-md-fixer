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

func TestSidecarPath_editedUsesOriginalSidecar(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	edited := filepath.Join(dir, "IMG_0378-edited.JPG")
	if err := os.WriteFile(edited, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	origJSON := filepath.Join(dir, "IMG_0378.JPG.json")
	if err := os.WriteFile(origJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(edited); got != origJSON {
		t.Fatalf("got %q want %q", got, origJSON)
	}
}

func TestSidecarPath_editedOwnSidecarTakesPrecedence(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	edited := filepath.Join(dir, "IMG_0378-edited.JPG")
	if err := os.WriteFile(edited, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	editedJSON := filepath.Join(dir, "IMG_0378-edited.JPG.json")
	if err := os.WriteFile(editedJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	origJSON := filepath.Join(dir, "IMG_0378.JPG.json")
	if err := os.WriteFile(origJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(edited); got != editedJSON {
		t.Fatalf("got %q want %q", got, editedJSON)
	}
}

func TestSidecarCleanupPaths_includesSupplementalSibling(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	media := filepath.Join(dir, "IMG_0452.JPG")
	jpgJSON := filepath.Join(dir, "IMG_0452.JPG.json")
	sup := filepath.Join(dir, "IMG_0452.JPG.supplemental-metadata.json")
	paths := SidecarCleanupPaths(media, jpgJSON)
	var hasJpg, hasSup bool
	for _, p := range paths {
		if filepath.Clean(p) == filepath.Clean(jpgJSON) {
			hasJpg = true
		}
		if filepath.Clean(p) == filepath.Clean(sup) {
			hasSup = true
		}
	}
	if !hasJpg || !hasSup {
		t.Fatalf("cleanup paths should include both .jpg.json and supplemental: %#v", paths)
	}
}

func TestSidecarPath_livePhotoMP4UsesJpgSupplemental(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mp4 := filepath.Join(dir, "IMG_5175.MP4")
	jpg := filepath.Join(dir, "IMG_5175.JPG")
	sup := filepath.Join(dir, "IMG_5175.JPG.supplemental-metadata.json")
	for _, p := range []struct {
		path string
		data string
	}{
		{mp4, "m"},
		{jpg, "j"},
		{sup, "{}"},
	} {
		if err := os.WriteFile(p.path, []byte(p.data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got := SidecarPath(mp4); got != sup {
		t.Fatalf("got %q want %q", got, sup)
	}
}

func TestSidecarPath_livePhotoMP4UsesJpgJson(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mp4 := filepath.Join(dir, "clip.MP4")
	jpg := filepath.Join(dir, "clip.JPG")
	jpgJSON := filepath.Join(dir, "clip.JPG.json")
	for _, p := range []string{mp4, jpg} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(jpgJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(mp4); got != jpgJSON {
		t.Fatalf("got %q want %q", got, jpgJSON)
	}
}

func TestSidecarPath_motionOwnSidecarBeforePairedStill(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mov := filepath.Join(dir, "IMG_5134.MOV")
	jpg := filepath.Join(dir, "IMG_5134.JPG")
	movSup := filepath.Join(dir, "IMG_5134.MOV.supplemental-metadata.json")
	jpgSup := filepath.Join(dir, "IMG_5134.JPG.supplemental-metadata.json")
	for _, p := range []string{mov, jpg} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(movSup, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(jpgSup, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(mov); got != movSup {
		t.Fatalf("got %q want %q (motion file should use its own sidecar first)", got, movSup)
	}
}

func TestSidecarBorrowedFromDifferentMedia_livePhoto(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mp4 := filepath.Join(dir, "IMG_5175.MP4")
	sup := filepath.Join(dir, "IMG_5175.JPG.supplemental-metadata.json")
	if !SidecarBorrowedFromDifferentMedia(mp4, sup) {
		t.Fatal("expected borrowed for MP4 using JPG sidecar")
	}
	ownSup := mp4 + ".supplemental-metadata.json"
	if SidecarBorrowedFromDifferentMedia(mp4, ownSup) {
		t.Fatal("same-basename sidecar is not borrowed")
	}
	jpg := filepath.Join(dir, "IMG_5175.JPG")
	if SidecarBorrowedFromDifferentMedia(jpg, sup) {
		t.Fatal("still file using its own sidecar is not borrowed")
	}
}

func TestSidecarDeletionPaths_stillWithMotionSiblingDeletesNothing(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	jpg := filepath.Join(dir, "IMG_5175.JPG")
	mp4 := filepath.Join(dir, "IMG_5175.MP4")
	sup := filepath.Join(dir, "IMG_5175.JPG.supplemental-metadata.json")
	for _, p := range []string{jpg, mp4} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(sup, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	paths := SidecarDeletionPaths(jpg, sup)
	if len(paths) != 0 {
		t.Fatalf("expected no deletions for still when motion sibling exists (same-run MP4 needs JSON), got %#v", paths)
	}
}

func TestSidecarDeletionPaths_omitsBorrowedResolvedJson(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mp4 := filepath.Join(dir, "IMG_5175.MP4")
	sup := filepath.Join(dir, "IMG_5175.JPG.supplemental-metadata.json")
	if err := os.WriteFile(mp4, []byte("m"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sup, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	paths := SidecarDeletionPaths(mp4, sup)
	for _, p := range paths {
		if filepath.Clean(p) == filepath.Clean(sup) {
			t.Fatalf("borrowed sidecar should not be deleted: still in %#v", paths)
		}
	}
}

func TestSidecarPath_editedCaseInsensitiveSuffix(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	edited := filepath.Join(dir, "IMG_0378-EDITED.jpg")
	if err := os.WriteFile(edited, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	origJSON := filepath.Join(dir, "IMG_0378.jpg.json")
	if err := os.WriteFile(origJSON, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := SidecarPath(edited); got != origJSON {
		t.Fatalf("got %q want %q", got, origJSON)
	}
}
