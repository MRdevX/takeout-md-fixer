package takeout

// MediaExtensions lists file extensions treated as photos/videos from Takeout.
var MediaExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".heic": true, ".webp": true, ".bmp": true, ".tiff": true,
	".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	".3gp": true, ".m4v": true,
}
