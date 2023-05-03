package file

import (
	"bytes"
	"io"
	"os"

	"github.com/ipfs/boxo/files"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

type Reader struct {
	*files.ReaderFile
}

func (r *Reader) GetBlock(offset, size int64, ci cid.Cid) (blocks.Block, error) {
	rawBlock := make([]byte, size)

	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	data := make([]byte, 50)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}

	block, err := blocks.NewBlockWithCid(rawBlock, ci)
	if err != nil {
		return nil, err
	}
	return block, nil
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

func NewWriterFromBytes(bt []byte) *Writer {
	writer := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return NewWriter(writer)
}

func NewWriter(node files.Node) *Writer {
	return &Writer{
		Node: node,
	}
}
