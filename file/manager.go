package file

import (
	"context"
	"path/filepath"

	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
)

// file manager by paths or cids
type FileManager struct {
	cids map[cid.Cid]*FileInfo // map[cids]FileInfo
}

func NewFileManager() *FileManager {
	return &FileManager{
		cids: map[cid.Cid]*FileInfo{},
	}
}

func (fm *FileManager) PutNode(nd format.Node, path string, blockSize int) {
	size, _ := nd.Size()
	fm.Put(nd.Cid(), path, 0, int(size)) // put root cid
	for i, link := range nd.Links() {    // put link cids
		fm.Put(link.Cid, filepath.Join(path, link.Name), i*blockSize, blockSize)
	}
}

func (fm *FileManager) Put(ci cid.Cid, path string, offset, size int) {
	// add cids
	fm.cids[ci] = NewFileInfo(ci, path, offset, size)
}

func (fm *FileManager) Get(cid cid.Cid) *FileInfo {
	return fm.cids[cid]
}

func (f *FileManager) Has(ctx context.Context, ci cid.Cid) bool {
	_, ok := f.cids[ci]
	return ok
}

func (f *FileManager) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	out := make(chan cid.Cid)
	for ci := range f.cids {
		out <- ci
	}
	return out, nil
}

type FileInfo struct {
	Path   string
	Ci     cid.Cid
	Offset int
	Size   int
}

func NewFileInfo(ci cid.Cid, path string, offset, size int) *FileInfo {
	return &FileInfo{
		Path:   path,
		Ci:     ci,
		Offset: offset,
		Size:   size,
	}
}
