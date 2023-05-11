package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/boxo/files"
)

// Store is a file store that provides methods for reading, writing, and deleting files.
type Store struct {
	// rootPath is the root directory of the file store.
	rootPath string
}

// NewFileStore creates a new file store with the given root directory.
func NewFileStore(rootPath string) *Store {
	os.MkdirAll(rootPath, 0755)

	return &Store{
		rootPath: rootPath,
	}
}

// Put writes the contents of the given writer to the file at the given path.
func (f *Store) Put(ctx context.Context, path string, writer *Writer) error {
	fileName := filepath.Join(f.rootPath, path)
	filePath := filepath.Dir(fileName)
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return err
	}

	if writer != nil {
		err = files.WriteTo(writer.Node, fileName)
		if err != nil {
			return err
		}
		err = writer.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Get returns a reader for the file at the given path.
func (f *Store) Get(ctx context.Context, path string) (*Reader, error) {
	fullPath := filepath.Join(f.rootPath, path)
	return NewReaderFromPath(fullPath), nil
}

// Delete deletes the file at the given path.
func (f *Store) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(f.rootPath, path)
	return os.RemoveAll(fullPath)
}

// Utility functions

// Overwrite overwrites the file at the given path with the contents of the given writer.
func (f *Store) Overwrite(ctx context.Context, path string, writer *Writer) error {
	err := f.Delete(ctx, path)
	if err != nil {
		return err
	}

	return f.Put(ctx, path, writer)
}

// Iterate iterates over the files in the directory at the given path, calling the given function for each file.
func (f *Store) Iterate(path string, fn func(fpath string, reader *Reader) error) error {
	fullPath := filepath.Join(f.rootPath, path)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	sf, err := files.NewSerialFile(fullPath, true, stat)
	if err != nil {
		return err
	}
	return files.Walk(sf, func(fpath string, node files.Node) error {
		return fn(fpath, NewReader(node))
	})
}
