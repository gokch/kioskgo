package file

import (
	"bytes"
	"io"
	"os"

	"github.com/ipfs/boxo/files"
)

type Reader struct {
	*files.ReaderFile
}

func (r *Reader) GetBlock(offset, size int) ([]byte, error) {
	if _, err := r.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, err
	}
	rawBlock := make([]byte, size)
	if _, err := r.Read(rawBlock); err != nil {
		return nil, err
	}

	return rawBlock, nil
}

func NewReaderFromPath(path string) *Reader {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	stat, err := file.Stat()
	if err != nil {
		return nil
	}

	reader, err := files.NewReaderPathFile(path, file, stat)
	if err != nil {
		return nil
	}

	return NewReader(reader)
}

func NewReader(reader *files.ReaderFile) *Reader {
	return &Reader{
		ReaderFile: reader,
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
