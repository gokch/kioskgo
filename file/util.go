package file

import (
	"os"
	"path/filepath"

	ds "github.com/ipfs/go-datastore"
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

func getFilename(rootPath string, path ds.Key) string {
	return filepath.Join(rootPath, path.String())
	// return filepath.Join(rootPath, path.String(), DEF_PATH_KEY)
}

// build ipld Node
