package takeout

import (
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Takeout uses ~46-rune truncation for very long basenames (filename limits).
const (
	takeoutBasenameTruncateRunes = 46
	longBasenameRunes            = 47
)

// motionPairVideoExts lists video extensions that may appear as the motion half of a
// Google Photos Live Photo–style pair; metadata often lives only on the still file.
var motionPairVideoExts = map[string]struct{}{
	".mp4": {}, ".mov": {}, ".avi": {}, ".mkv": {}, ".m4v": {}, ".3gp": {},
}

func isMotionPairVideoExt(ext string) bool {
	_, ok := motionPairVideoExts[ext]
	return ok
}

// pairedStillExtensions is the order we probe for a sibling still next to a motion file.
// Uppercase variants are tried first (common on iOS / Google Takeout) so returned paths match on disk.
var pairedStillExtensions = []string{
	".JPG", ".JPEG", ".jpg", ".jpeg", ".HEIC", ".heic",
}

func prefixRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func statOK(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// SidecarPath returns the path to the JSON metadata file for a media file, or "" if none.
//
// Resolution order (first existing file wins):
//  1. mediaPath + ".json" — e.g. photo.jpg.json
//  2. mediaPath + ".supplemental-metadata.json"
//  3. dir + nameNoExt + ".json" — e.g. photo.json (extension replaced; community Takeout variant)
//  4. If basename ≥47 runes: dir + first 46 runes of basename + ".json" (truncation variant)
//  5. Duplicate "name (n).ext" / "name(n).ext":
//     a. dir + stem + ".json" — e.g. IMG(1).jpg → IMG.json
//     b. dir + stem + ext + "(n)" + ".json" — e.g. IMG.jpg(1).json
//  6. Live Photo–style motion file: same stem as a still (e.g. IMG_5175.MP4 uses IMG_5175.JPG.supplemental-metadata.json)
//  7. Directory scan: truncated supplemental sidecars (*.supplemental*.json)
//  8. Google Photos edited export: IMG_0378-edited.JPG uses the same metadata as IMG_0378.JPG
//     (fallback only when no sidecar exists for the edited filename itself).
func SidecarPath(mediaPath string) string {
	dir := filepath.Dir(mediaPath)
	base := filepath.Base(mediaPath)
	ext := filepath.Ext(base)
	nameNoExt := strings.TrimSuffix(base, ext)

	candidates := []string{
		mediaPath + ".json",
		mediaPath + ".supplemental-metadata.json",
		filepath.Join(dir, nameNoExt+".json"),
	}

	if utf8.RuneCountInString(base) >= longBasenameRunes {
		candidates = append(candidates, filepath.Join(dir, prefixRunes(base, takeoutBasenameTruncateRunes)+".json"))
	}

	for _, c := range candidates {
		if statOK(c) {
			return c
		}
	}

	// Duplicate naming: IMG(1).jpg → IMG.json, else IMG.jpg(1).json
	if idx := strings.LastIndex(nameNoExt, "("); idx > 0 {
		suffix := nameNoExt[idx:]
		if strings.HasSuffix(suffix, ")") {
			originalName := nameNoExt[:idx]
			for _, c := range []string{
				filepath.Join(dir, originalName+".json"),
				filepath.Join(dir, originalName+ext+suffix+".json"),
			} {
				if statOK(c) {
					return c
				}
			}
		}
	}

	if p := sidecarFromPairedStill(dir, nameNoExt, ext); p != "" {
		return p
	}

	if p := sidecarSupplementalFromDir(dir, base); p != "" {
		return p
	}

	if stem, ok := editedOriginalStem(nameNoExt); ok {
		origPath := filepath.Join(dir, stem+ext)
		if origPath != mediaPath {
			if p := SidecarPath(origPath); p != "" {
				return p
			}
		}
	}

	return ""
}

// sidecarFromPairedStill resolves JSON that lives next to a sibling still image (Google Takeout Live Photo exports).
func sidecarFromPairedStill(dir, nameNoExt, mediaExt string) string {
	if !isMotionPairVideoExt(strings.ToLower(mediaExt)) {
		return ""
	}
	for _, stillExt := range pairedStillExtensions {
		stillBase := nameNoExt + stillExt
		stillPath := filepath.Join(dir, stillBase)
		for _, c := range []string{
			stillPath + ".json",
			stillPath + ".supplemental-metadata.json",
		} {
			if statOK(c) {
				return c
			}
		}
		if utf8.RuneCountInString(stillBase) >= longBasenameRunes {
			trunc := filepath.Join(dir, prefixRunes(stillBase, takeoutBasenameTruncateRunes)+".json")
			if statOK(trunc) {
				return trunc
			}
		}
	}
	return ""
}

// sidecarJSONOwnerBaseName returns the media basename a Takeout JSON file belongs to.
func sidecarJSONOwnerBaseName(jsonBase string) string {
	switch {
	case strings.HasSuffix(jsonBase, ".supplemental-metadata.json"):
		return strings.TrimSuffix(jsonBase, ".supplemental-metadata.json")
	case strings.HasSuffix(strings.ToLower(jsonBase), ".json"):
		return strings.TrimSuffix(jsonBase, ".json")
	default:
		return ""
	}
}

// SidecarCleanupRepresentative picks a media path for SidecarCleanupPaths when several files share one JSON
// (e.g. Live Photo JPG + MP4). Prefer the file basename that owns the sidecar filename.
func SidecarCleanupRepresentative(mediaPaths []string, resolvedJsonPath string) string {
	if len(mediaPaths) == 0 {
		return ""
	}
	owner := sidecarJSONOwnerBaseName(filepath.Base(resolvedJsonPath))
	if owner == "" {
		return mediaPaths[0]
	}
	for _, p := range mediaPaths {
		if strings.EqualFold(filepath.Base(p), owner) {
			return p
		}
	}
	return mediaPaths[0]
}

// editedOriginalStem reports whether nameNoExt ends with "-edited" (case-insensitive) and returns the stem without it.
func editedOriginalStem(nameNoExt string) (string, bool) {
	const suf = "-edited"
	if len(nameNoExt) <= len(suf) {
		return "", false
	}
	i := len(nameNoExt) - len(suf)
	if !strings.EqualFold(nameNoExt[i:], suf) {
		return "", false
	}
	return nameNoExt[:i], true
}

// SidecarCleanupPaths lists every path that may contain Takeout metadata for this media file.
// Google often ships both IMG_0452.JPG.json and IMG_0452.JPG.supplemental-metadata.json; SidecarPath
// only returns the first match, so deletion must remove all siblings.
func SidecarCleanupPaths(mediaPath, resolvedJsonPath string) []string {
	dir := filepath.Dir(mediaPath)
	base := filepath.Base(mediaPath)
	ext := filepath.Ext(base)
	nameNoExt := strings.TrimSuffix(base, ext)

	seen := make(map[string]bool)
	var out []string
	add := func(p string) {
		if p == "" {
			return
		}
		p = filepath.Clean(p)
		if seen[p] {
			return
		}
		seen[p] = true
		out = append(out, p)
	}

	add(resolvedJsonPath)
	add(mediaPath + ".json")
	add(mediaPath + ".supplemental-metadata.json")
	add(filepath.Join(dir, nameNoExt+".json"))
	if utf8.RuneCountInString(base) >= longBasenameRunes {
		add(filepath.Join(dir, prefixRunes(base, takeoutBasenameTruncateRunes)+".json"))
	}
	if idx := strings.LastIndex(nameNoExt, "("); idx > 0 {
		suffix := nameNoExt[idx:]
		if strings.HasSuffix(suffix, ")") {
			originalName := nameNoExt[:idx]
			add(filepath.Join(dir, originalName+".json"))
			add(filepath.Join(dir, originalName+ext+suffix+".json"))
		}
	}
	if p := sidecarSupplementalFromDir(dir, base); p != "" {
		add(p)
	}
	if stem, ok := editedOriginalStem(nameNoExt); ok {
		origPath := filepath.Join(dir, stem+ext)
		if origPath != mediaPath {
			origBase := filepath.Base(origPath)
			origNameNoExt := strings.TrimSuffix(origBase, ext)
			add(origPath + ".json")
			add(origPath + ".supplemental-metadata.json")
			add(filepath.Join(dir, origNameNoExt+".json"))
			if p := sidecarSupplementalFromDir(dir, origBase); p != "" {
				add(p)
			}
		}
	}

	return out
}

// sidecarSupplementalFromDir matches truncated supplemental-metadata filenames, e.g.
// verylong….jpg.supplemental-metad.json (Google truncates the suffix after ~46 chars).
func sidecarSupplementalFromDir(dir, base string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	prefixFull := base + "."
	prefix46 := prefixRunes(base, takeoutBasenameTruncateRunes) + "."
	long := utf8.RuneCountInString(base) >= longBasenameRunes

	var matches []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}
		lower := strings.ToLower(name)
		if !strings.Contains(lower, "supplemental") {
			continue
		}
		if strings.HasPrefix(name, prefixFull) {
			matches = append(matches, filepath.Join(dir, name))
			continue
		}
		if long && strings.HasPrefix(name, prefix46) {
			matches = append(matches, filepath.Join(dir, name))
		}
	}
	if len(matches) == 0 {
		return ""
	}
	// Deterministic pick when multiple (rare)
	best := matches[0]
	for _, m := range matches[1:] {
		if m < best {
			best = m
		}
	}
	return best
}
