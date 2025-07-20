package utils

import (
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/api/drive/v3"
)

func CreateFolder(srv *drive.Service, directory string, parentId string) (*drive.File, error) {
	dir := path.Base(directory)
	var folder *drive.File
	if parentId != "" {
		folder = &drive.File{
			Name:     dir,
			MimeType: "application/vnd.google-apps.folder",
			Parents:  []string{parentId},
		}
	} else {
		folder = &drive.File{
			Name:     dir,
			MimeType: "application/vnd.google-apps.folder",
		}
	}
	f, err := srv.Files.Create(folder).Do()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetCurrentDirectoryMap(data map[string]interface{}, version string, wd string, fullPath string) (map[string]interface{}, bool) {
	wd = filepath.Clean(wd)
    fullPath = filepath.Clean(fullPath)
	newMap, ok := RecursiveGetMap(data, version, wd)
	if !ok {
		return nil, false
	}

	var parts []string
	if filepath.Clean(fullPath) != filepath.Clean(wd) {
		relative := strings.TrimPrefix(filepath.ToSlash(fullPath), filepath.ToSlash(wd)+"/")
		parts = strings.Split(relative, "/")

	} else {
		parts = []string{}
	}
	newMap = RecursiveCreateMap(newMap, parts...)
	return newMap, true
}
