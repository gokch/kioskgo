package file

import (
	"path/filepath"
	"runtime/debug"
	"sync"

	"github.com/dghubble/trie"
	"github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/go-cid"
)

// file manager by paths or cids
type FileManager struct {
	mtx sync.Mutex

	rootPath string

	cids  map[cid.Cid]*FileInfo // map[cids]FileInfo
	paths *trie.PathTrie        // trie[fileInfo]
}

func NewFileManager(rootPath string) *FileManager {
	return &FileManager{
		mtx:      sync.Mutex{},
		rootPath: rootPath,
		cids:     map[cid.Cid]*FileInfo{},
		paths:    trie.NewPathTrie(),
	}
}

func (fm *FileManager) Exist(path string, ci cid.Cid) bool {
	return fm.Get(path, ci) != nil
}

func (fm *FileManager) Put(path string, ci cid.Cid) {
	pathWithCid := filepath.Join(path, ci.String())

	fi := NewFileInfo(fm.rootPath)
	fi.Set(pathWithCid)

	// add paths
	fm.paths.Put(pathWithCid, fi)

	// add cids
	fm.cids[ci] = fi
}

func (fm *FileManager) Delete(path string, ci cid.Cid) {
	unixfs.NewFSNode()
	// del paths
	pathWithCid := filepath.Join(path, ci.String())
	fm.paths.Delete(pathWithCid)

	// del cids
	delete(fm.cids, ci)
}

func (fm *FileManager) Get(path string, ci cid.Cid) *FileInfo {
	pathWithCid := filepath.Join(path, ci.String())

	return fm.paths.Get(pathWithCid).(*FileInfo)
}

func (fm *FileManager) GetCid(cid cid.Cid) *FileInfo {
	return fm.cids[cid]
}

func (fm *FileManager) Walk(path string) []*FileInfo {
	fis := make([]*FileInfo, 0, 1024)
	fm.paths.WalkPath(path, func(key string, value interface{}) error {
		fis = append(fis, value.(*FileInfo))
		return nil
	})
	return fis
}

func (fm *FileManager) Clear() {
	fm.mtx.Lock()
	defer fm.mtx.Unlock()

	fm.cids = map[cid.Cid]*FileInfo{}
	fm.paths = trie.NewPathTrie()

	// clear orphan memory
	debug.FreeOSMemory()
}

type FileInfo struct {
	rootPath     string
	relativePath string
	ci           cid.Cid
}

func NewFileInfo(rootPath string) *FileInfo {
	return &FileInfo{
		rootPath: rootPath,
	}
}

func (fi *FileInfo) Set(path string) error {
	var err error
	// if full path, make relative path
	if filepath.HasPrefix(fi.rootPath, path) {
		path, err = filepath.Rel(fi.rootPath, path)
		if err != nil {
			return err
		}
	}

	// extract specific cid
	ci, err := cid.Parse(filepath.Base(path))
	if err != nil {
		return err
	}

	fi.relativePath = path
	fi.ci = ci
	return nil
}

func (fi *FileInfo) GetFullPath() (fullPath string) {
	return filepath.Join(fi.rootPath, fi.relativePath)
}
