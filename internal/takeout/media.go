package takeout

// Extensions treated as photos/videos from Takeout (lowercase keys).
var mediaExtensions = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {},
	".heic": {}, ".webp": {}, ".bmp": {}, ".tiff": {},
	".mp4": {}, ".mov": {}, ".avi": {}, ".mkv": {},
	".3gp": {}, ".m4v": {},
}

func isMediaExt(ext string) bool {
	_, ok := mediaExtensions[ext]
	return ok
}
