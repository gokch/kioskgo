package file

import (
	"bytes"
	"os"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
)

type Reader struct {
	*files.ReaderFile
	Cids []cid.Cid // cids in specific path, if withCid == nil, Cid is not specified
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
	return &Reader{ReaderFile: reader, Cids: nil}
}

type Writer struct {
	files.Node
	Cid []cid.Cid // cids in specific path, if withCid == nil, Cid is not specified
}

func NewWriterFromBytes(bt []byte) *Writer {
	writer := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return NewWriter(writer)
}

func NewWriter(node files.Node) *Writer {
	return &Writer{Node: node, Cid: nil}
}
