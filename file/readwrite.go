package file

import (
	"bytes"
	"os"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
)

type Reader struct {
	*files.ReaderFile
	Cid cid.Cid
}

func NewReaderFromPath(path string) *Reader {
	open, err := os.Open(path)
	if err != nil {
		return nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil
	}

	reader, err := files.NewReaderPathFile(path, open, stat)
	if err != nil {
		return nil
	}

	return NewReader(reader)
}

func NewReaderFromBytes(bt []byte) *Reader {
	reader := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return NewReader(reader)
}

func NewReader(reader *files.ReaderFile) *Reader {
	return &Reader{ReaderFile: reader, Cid: cid.Cid{}}
}

type Writer struct {
	files.Node
	Cid cid.Cid
}

func NewWriterFromBytes(bt []byte) *Writer {
	writer := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return NewWriter(writer)
}

func NewWriter(node files.Node) *Writer {
	return &Writer{Node: node, Cid: cid.Cid{}}
}
