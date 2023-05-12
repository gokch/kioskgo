package file

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/ipfs/boxo/files"
)

type Reader struct {
	files.Node
}

func (r *Reader) Get() ([]byte, error) {
	return r.GetBlock(0, 0)
}

func (r *Reader) GetBlock(offset, size int64) ([]byte, error) {
	var err error

	file, ok := r.Node.(*files.ReaderFile)
	if ok != true {
		return nil, errors.New("cannot read block from directory")
	}

	if offset < 0 {
		offset = 0
	}
	if size <= 0 {
		size, err = file.Size()
		if err != nil {
			return nil, err
		}
	}

	if _, err := file.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, err
	}
	rawBlock := make([]byte, size)
	if _, err := file.Read(rawBlock); err != nil {
		return nil, err
	}

	return rawBlock, nil
}

func NewReaderFromPath(path string) *Reader {
	stat, err := os.Stat(path)
	if err != nil {
		return nil
	}

	reader, err := files.NewSerialFile(path, true, stat)
	if err != nil {
		return nil
	}

	return NewReader(reader)
}

func NewReader(reader files.Node) *Reader {
	return &Reader{
		Node: reader,
	}
}

type Writer struct {
	files.Node
}

func NewWriterFromPath(path string) *Writer {
	var err error

	stat, err := os.Stat(path)
	if err != nil {
		return nil
	}

	nd, err := files.NewSerialFile(path, true, stat)
	if err != nil {
		return nil
	}
	return NewWriter(nd)
}

func NewWriterFromBytes(bt []byte) *Writer {
	writer := files.NewReaderFile(bytes.NewReader(bt))
	return NewWriter(writer)
}

func NewWriter(node files.Node) *Writer {
	return &Writer{
		Node: node,
	}
}
