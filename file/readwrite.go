package file

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
)

type Reader struct {
	*files.ReaderFile
	ci cid.Cid
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

func NewReader(reader *files.ReaderFile) *Reader {
	// get cid
	var ci cid.Cid
	cidInfoPath := filepath.Join(filepath.Dir(reader.AbsPath()), DEF_PATH_CID_INFO)
	rawCid, err := os.ReadFile(cidInfoPath)
	if err == nil {
		_, c, err := cid.CidFromBytes(rawCid)
		if err == nil {
			ci = c
		}
	}

	return &Reader{
		ReaderFile: reader,
		ci:         ci,
	}
}

type Writer struct {
	files.Node
	ci cid.Cid
}

func NewWriterFromBytes(bt []byte, ci cid.Cid) *Writer {
	writer := files.NewReaderFile(bytes.NewReader(bt)).(*files.ReaderFile)
	return NewWriter(writer, ci)
}

func NewWriter(node files.Node, ci cid.Cid) *Writer {
	return &Writer{
		Node: node,
		ci:   ci,
	}
}

// cids
const (
	DEF_PATH_CID_INFO = "/.info"
)
