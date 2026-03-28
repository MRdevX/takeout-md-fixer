package takeout

// TakeoutMeta matches Google Takeout JSON sidecar fields used for EXIF/GPS.
type TakeoutMeta struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	PhotoTakenTime TimeStamp `json:"photoTakenTime"`
	GeoData        GeoData   `json:"geoData"`
	GeoDataExif    GeoData   `json:"geoDataExif"`
}

type TimeStamp struct {
	Timestamp string `json:"timestamp"`
	Formatted string `json:"formatted"`
}

type GeoData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type MediaFile struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	JsonPath string `json:"jsonPath"`
	HasJson  bool   `json:"hasJson"`
	Status   string `json:"status"`
	Error    string `json:"error"`
}

type ScanResult struct {
	FolderPath  string      `json:"folderPath"`
	TotalMedia  int         `json:"totalMedia"`
	WithJson    int         `json:"withJson"`
	WithoutJson int         `json:"withoutJson"`
	Files       []MediaFile `json:"files"`
}

type FixProgress struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	File    string `json:"file"`
	Status  string `json:"status"`
}

type FixResult struct {
	Total            int `json:"total"`
	Success          int `json:"success"`
	Failed           int `json:"failed"`
	Skipped          int `json:"skipped"`
	JsonDeleted      int `json:"jsonDeleted"`
	JsonDeleteFailed int `json:"jsonDeleteFailed"`
}
