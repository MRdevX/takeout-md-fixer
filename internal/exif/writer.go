package exif

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	exiftool "github.com/barasher/go-exiftool"

	"takeout-md-fixer/internal/takeout"
)

// Writer applies Takeout metadata to media files via ExifTool.
type Writer struct {
	et *exiftool.Exiftool
}

// NewWriter starts an ExifTool session. Caller must Close when done.
func NewWriter() (*Writer, error) {
	path, err := ResolveExiftoolPath()
	if err != nil {
		return nil, err
	}
	et, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(path))
	if err != nil {
		return nil, fmt.Errorf("failed to start ExifTool at %s: %w", path, err)
	}
	return &Writer{et: et}, nil
}

// Close releases ExifTool resources.
func (w *Writer) Close() {
	if w.et != nil {
		w.et.Close()
	}
}

// WriteMetadata writes EXIF/GPS and related tags from meta into mediaPath.
func (w *Writer) WriteMetadata(mediaPath string, meta *takeout.TakeoutMeta) error {
	fi := exiftool.FileMetadata{
		File:   mediaPath,
		Fields: make(map[string]interface{}),
	}

	var photoTime time.Time
	var hasPhotoTime bool
	if meta.PhotoTakenTime.Timestamp != "" {
		ts, err := strconv.ParseInt(meta.PhotoTakenTime.Timestamp, 10, 64)
		if err == nil {
			photoTime = time.Unix(ts, 0).UTC()
			hasPhotoTime = true
			dateStr := photoTime.Format("2006:01:02 15:04:05")
			fi.SetString("DateTimeOriginal", dateStr)
			fi.SetString("CreateDate", dateStr)
			fi.SetString("ModifyDate", dateStr)
			fi.SetString("FileModifyDate", dateStr)
			fi.SetString("FileCreateDate", dateStr)
		}
	}

	if meta.GeoData.Latitude != 0 || meta.GeoData.Longitude != 0 {
		fi.SetFloat("GPSLatitude", meta.GeoData.Latitude)
		fi.SetFloat("GPSLongitude", meta.GeoData.Longitude)

		latRef := "N"
		if meta.GeoData.Latitude < 0 {
			latRef = "S"
		}
		lngRef := "E"
		if meta.GeoData.Longitude < 0 {
			lngRef = "W"
		}
		fi.SetString("GPSLatitudeRef", latRef)
		fi.SetString("GPSLongitudeRef", lngRef)

		if meta.GeoData.Altitude != 0 {
			fi.SetFloat("GPSAltitude", meta.GeoData.Altitude)
		}
	}

	if meta.Description != "" {
		fi.SetString("ImageDescription", meta.Description)
	}

	ext := strings.ToLower(filepath.Ext(mediaPath))
	if hasPhotoTime && (ext == ".mp4" || ext == ".mov" || ext == ".m4v" || ext == ".3gp") {
		dateStr := photoTime.Format("2006:01:02 15:04:05")
		fi.SetString("MediaCreateDate", dateStr)
		fi.SetString("MediaModifyDate", dateStr)
		fi.SetString("TrackCreateDate", dateStr)
		fi.SetString("TrackModifyDate", dateStr)
	}

	batch := []exiftool.FileMetadata{fi}
	w.et.WriteMetadata(batch)
	if batch[0].Err != nil {
		return fmt.Errorf("error writing metadata to %s: %w", mediaPath, batch[0].Err)
	}

	if hasPhotoTime {
		if err := os.Chtimes(mediaPath, photoTime, photoTime); err != nil {
			return fmt.Errorf("metadata written but could not set file timestamps: %w", err)
		}
	}
	return nil
}
