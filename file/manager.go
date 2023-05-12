package file

import (
	"context"
	"path/filepath"

	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
)

// Manager is a file manager that provides methods for reading, writing, and deleting files.
type Manager struct {
	// cidMapper is a map of CIDs to file information.
	cidMapper map[cid.Cid]*FileInfo
}

// NewFileManager creates a new file manager.
func NewFileManager() *Manager {
	return &Manager{
		cidMapper: map[cid.Cid]*FileInfo{},
	}
}

// PutNode writes the given node to the file manager.
func (fm *Manager) PutNode(nd format.Node, path string, blockSize int64) {
	size, _ := nd.Size()
	fm.Put(nd.Cid(), path, 0, int64(size)) // put root cid
	for i, link := range nd.Links() {      // put link cidMapper
		fm.Put(link.Cid, filepath.Join(path, link.Name), int64(i)*blockSize, blockSize)
	}
}

// Put writes the given CID and file information to the file manager.
func (fm *Manager) Put(ci cid.Cid, path string, offset, size int64) {
	// add cidMapper
	fm.cidMapper[ci] = NewFileInfo(ci, path, offset, size)
}

// Get returns the file information for the given CID.
func (fm *Manager) Get(cid cid.Cid) *FileInfo {
	return fm.cidMapper[cid]
}

// Has checks if the file manager contains the given CID.
func (fm *Manager) Has(ctx context.Context, ci cid.Cid) bool {
	_, ok := fm.cidMapper[ci]
	return ok
}

// AllKeysChan returns a channel that emits all of the CIDs in the file manager.
func (fm *Manager) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	out := make(chan cid.Cid)
	for ci := range fm.cidMapper {
		out <- ci
	}
	return out, nil
}

// FileInfo is a struct that contains information about a file.
type FileInfo struct {
	// Path is the path to the file.
	Path string
	// Ci is the CID of the file.
	Ci cid.Cid
	// Offset is the offset of the file in the file system.
	Offset int64
	// Size is the size of the file.
	Size int64
}

// NewFileInfo creates a new FileInfo struct.
func NewFileInfo(ci cid.Cid, path string, offset, size int64) *FileInfo {
	return &FileInfo{
		Path:   path,
		Ci:     ci,
		Offset: offset,
		Size:   size,
	}
}
