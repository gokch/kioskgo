package file

import (
	"runtime/debug"
	"sync"

	"github.com/dghubble/trie"
	"github.com/ipfs/go-cid"
)

// file manager by paths or cids
type FileManager struct {
	mtx sync.Mutex

	cids  map[cid.Cid]string
	paths *trie.PathTrie
}

func NewFileManager() *FileManager {
	return &FileManager{
		cids:  map[cid.Cid]string{},
		paths: trie.NewPathTrie(),
	}
}

func (fm *FileManager) GetCids(cid cid.Cid) string {
	return fm.cids[cid]
}

func (fm *FileManager) Put(path string, fi *fileInfo) {
	fm.paths.Put(path, fi)
}

func (fm *FileManager) Delete(path string) {
	fm.paths.Delete(path)
}

func (fm *FileManager) Get(path string) *fileInfo {
	return fm.paths.Get(path).(*fileInfo)
}

func (fm *FileManager) Walk(path string) []*fileInfo {
	fis := make([]*fileInfo, 0, 1024)
	fm.paths.WalkPath(path, func(key string, value interface{}) error {
		fis = append(fis, value.(*fileInfo))
		return nil
	})
	return fis
}

func (fm *FileManager) Clear() {
	fm.mtx.Lock()
	defer fm.mtx.Unlock()

	fm.cids = map[cid.Cid]string{}
	fm.paths = trie.NewPathTrie()

	// clear orphan memory
	debug.FreeOSMemory()
}

type fileInfo struct {
	size uint64

	rootPath   string
	parentPath string
	myPath     string // cids? TODO
	childsPath []string
}
