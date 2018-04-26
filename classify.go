package backup

import (
	"path/filepath"
	"strings"
)

const (
	unclassifiedFile = iota
	imageFile
	videoFile
)

// TODO: This is currently conservative. A more robust approach using file(1) or
// `http.DetectContentType` would be nice, although both are hacky.
//
// Also see https://en.wikipedia.org/wiki/Raw_image_format#File_contents
var imageExtensions = map[string]bool{
	".arw":  true,
	".bmp":  true,
	".cr2":  true,
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
}

// classifyExt classifies `ext`, expecting a leading period. `ext` will be
// normalized to lowercase first.
func classifyExt(ext string) int {
	extLower := strings.ToLower(ext)
	switch {
	case imageExtensions[extLower]:
		return imageFile
	case videoExtensions[extLower]:
		return videoFile
	default:
		return unclassifiedFile
	}
}

func classifyPath(path string) int {
	return classifyExt(filepath.Ext(path))
}
