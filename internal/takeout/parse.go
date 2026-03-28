package takeout

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseMetadataFile reads and unmarshals a Takeout JSON sidecar.
func ParseMetadataFile(jsonPath string) (*TakeoutMeta, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", jsonPath, err)
	}

	var meta TakeoutMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", jsonPath, err)
	}

	if meta.GeoDataExif.Latitude != 0 || meta.GeoDataExif.Longitude != 0 {
		meta.GeoData = meta.GeoDataExif
	}

	return &meta, nil
}
