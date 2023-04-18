package file

import (
	"path/filepath"
	"runtime/debug"
	"sync"

	"github.com/dghubble/trie"
	"github.com/ipfs/go-cid"
)

// file manager by paths or cids
type FileManager struct {
	mtx sync.Mutex

	rootPath string

	cids  map[cid.Cid]*FileInfo
	paths *trie.PathTrie
}

func NewFileManager(rootPath string) *FileManager {
	return &FileManager{
		rootPath: rootPath,
		cids:     map[cid.Cid]*FileInfo{},
		paths:    trie.NewPathTrie(),
	}
}

func (fm *FileManager) GetCids(cid cid.Cid) *FileInfo {
	return fm.cids[cid]
}

func (fm *FileManager) Exist(path string, ci cid.Cid) bool {
	return fm.Get(path, ci) != nil
}

func (fm *FileManager) Put(path string, ci cid.Cid) {
	fi := &FileInfo{
		rootPath: fm.rootPath,
		myPath:   path,
		cidPath:  ci.String(),
	}

	// add paths
	pathWithCid := filepath.Join(path, ci.String())
	fm.paths.Put(pathWithCid, fi)

	// add cids
	fm.cids[ci] = fi
}

func (fm *FileManager) Delete(path string, ci cid.Cid) {
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

func (fm *FileManager) GetByCid(cid cid.Cid) *FileInfo {
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
	rootPath string
	myPath   string
	cidPath  string
}
