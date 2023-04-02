package file

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func NewFileSystem(root string) *FileSystem {
	fs := &FileSystem{}
	err := fs.Init(root)
	if err != nil {
		// TODO : add logger
		return nil
	}
	return fs
}

type FileSystem struct {
	root string
	iofs afero.IOFS

	bufsize int
}

func (f *FileSystem) Init(root string) error {
	f.root = root
	f.iofs = afero.NewIOFS(afero.NewOsFs())
	err := f.iofs.MkdirAll(f.root, 0755)
	if err != nil {
		return err
	}
	f.bufsize = 1024 // 1kb
	return nil
}

func (f *FileSystem) Add(path string, name string, reader io.Reader) error {
	fullpath := filepath.Join(f.root, path)
	filepath := filepath.Join(fullpath, name)

	err := f.iofs.MkdirAll(fullpath, 0755)
	if err != nil {
		return err
	}

	return afero.WriteReader(f.iofs.Fs, filepath, reader)
}

func (f *FileSystem) Get(path string, name string) (io.Reader, error) {
	fullpath := filepath.Join(f.root, path)
	filepath := filepath.Join(fullpath, name)

	return f.iofs.Fs.Open(filepath)
}

func (f *FileSystem) Iterate(path string, fn func(reader io.Reader) error) error {
	fullpath := filepath.Join(f.root, path)

	return afero.Walk(f.iofs.Fs, fullpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := f.iofs.Fs.Open(path)
			if err != nil {
				return nil
			}
			return fn(file)
		}
		return nil
	})
}
