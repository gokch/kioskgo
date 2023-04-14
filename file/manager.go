package file

import "github.com/dghubble/trie"

// manager
type fileManager struct {
	t *trie.PathTrie
}

func NewPaths() *fileManager {
	return &fileManager{t: trie.NewPathTrie()}
}

func (p *fileManager) Add(path string) {
	p.t.Put(path, true)
}
