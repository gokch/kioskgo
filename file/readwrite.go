package file

import (
	"os"
)

// cids
const (
	DEF_PATH_CID_INFO = "/.info"
	DEF_PATH_KEY      = "/.key"
)

// isDir returns whether given path is a directory
func isDir(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return finfo.IsDir()
}

// isFile returns whether given path is a file
func isFile(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !finfo.IsDir()
}
