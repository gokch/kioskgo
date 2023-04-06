package file

import (
	"bytes"

	"github.com/ipfs/boxo/files"
)

type Reader struct {
	*files.ReaderFile
}

func NewReaderFromBytes(bt []byte) *Reader {
	reader := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return &Reader{reader}
}

func NewReader(reader *files.ReaderFile) *Reader {
	return &Reader{reader}
}

type Writer struct {
	files.Node
}

func NewWriter(node files.Node) *Writer {
	return &Writer{node}
}

func NewWriterFromBytes(bt []byte) *Writer {
	writer := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return &Writer{writer}
}
