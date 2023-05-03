package file

import (
	"github.com/ipfs/go-cid"
)

// file manager by paths or cids
type FileManager struct {
	blockSize uint64
	cids      map[cid.Cid]*FileInfo // map[cids]FileInfo
	// paths     *trie.PathTrie        // trie[fileInfo]
}

func NewFileManager(blockSize uint64) *FileManager {
	return &FileManager{
		blockSize: blockSize,
		cids:      map[cid.Cid]*FileInfo{},
	}
}

func (fm *FileManager) Put(ci cid.Cid, path string, offset uint64) {
	// add cids
	fm.cids[ci] = NewFileInfo(ci, path, offset, fm.blockSize)
}

func (fm *FileManager) Get(cid cid.Cid) *FileInfo {
	return fm.cids[cid]
}

type FileInfo struct {
	path   string
	ci     cid.Cid
	offset uint64
	size   uint64
}

func NewFileInfo(ci cid.Cid, path string, offset, size uint64) *FileInfo {
	return &FileInfo{
		path:   path,
		ci:     ci,
		offset: offset,
		size:   size,
	}
}
