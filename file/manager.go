package file

import (
	"github.com/dghubble/trie"
	"github.com/ipfs/go-cid"
)

// file manage by paths
type fileManager struct {
	paths *trie.PathTrie
}

func NewManager() *fileManager {
	return &fileManager{paths: trie.NewPathTrie()}
}

func (fm *fileManager) Put(path string, fi *fileInfo) {
	fm.paths.Put(path, fi)
}

func (fm *fileManager) Delete(path string) {
	fm.paths.Delete(path)
}

type fileInfo struct {
	path string // without cid path
	cid  cid.Cid
}
