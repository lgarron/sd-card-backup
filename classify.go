package backup

import (
	"path/filepath"
	"strings"
)

type fileClassification int

const (
	unclassifiedFile = iota
	imageFile
	videoFile
	rawVideoFile
	audioFile
)

// TODO: This is currently conservative. A more robust approach using file(1) or
// `http.DetectContentType` would be nice, although both are hacky.
//
// Also see https://en.wikipedia.org/wiki/Raw_image_format#File_contents
var imageExtensions = map[string]bool{
	".arw":  true,
	".bmp":  true,
	".cr2":  true,
	".cr3":  true,
	".dng":  true,
	".gif":  true,
	".jpeg": true,
	".jpg":  true,
	".nef":  true,
	".png":  true,
	".raw":  true,
	".tif":  true,
	".webm": true,
}

// TODO: This is currently conservative. A more robust approach using file(1) or
// `http.DetectContentType` would be nice, although both are hacky.
var videoExtensions = map[string]bool{
	".avi":   true,
	".m4v":   true,
	".mkv":   true,
	".mov":   true,
	".mp4":   true,
	".mpeg:": true,
	".mpg:":  true,
	".mts":   true,
	".mxf":   true,
}

// TODO: This is currently conservative. A more robust approach using file(1) or
// `http.DetectContentType` would be nice, although both are hacky.
var rawVideoExtensions = map[string]bool{
	".crm": true, // Canon raw movie
}

// TODO: This is currently conservative. A more robust approach using file(1) or
// `http.DetectContentType` would be nice, although both are hacky.
var audioExtensions = map[string]bool{
	".aac":  true,
	".aif":  true,
	".aiff": true,
	".flac": true,
	".m4a":  true,
	".mp3":  true,
	".ogg":  true,
	".wav":  true,
	".wma":  true,
}

// classifyExt classifies `ext`, expecting a leading period. `ext` will be
// normalized to lowercase first.
func classifyExt(ext string) fileClassification {
	extLower := strings.ToLower(ext)
	switch {
	case imageExtensions[extLower]:
		return imageFile
	case videoExtensions[extLower]:
		return videoFile
	case rawVideoExtensions[extLower]:
		return rawVideoFile
	case audioExtensions[extLower]:
		return audioFile
	default:
		return unclassifiedFile
	}
}

func classifyPath(path string) fileClassification {
	return classifyExt(filepath.Ext(path))
}
