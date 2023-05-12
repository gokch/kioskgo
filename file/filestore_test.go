package file

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/ipfs/boxo/files"
	"github.com/stretchr/testify/require"
)

func TestStoreNewGet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test")
	err := fs.Put(ctx, "test/abc/d/e.txt", NewWriterFromBytes(data1))
	require.NoError(t, err)

	reader, err := fs.Get(ctx, "test/abc/d/e.txt")
	require.NoError(t, err)

	data2, err := reader.Get()
	require.NoError(t, err)
	reader.Close()

	require.Equal(t, data1, data2)

	err = fs.Delete(ctx, "")
	require.NoError(t, err)
}

func TestStoreFolder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("testtest")
	err := fs.Put(ctx, "test/abc/d/e.txt", NewWriterFromBytes(data1))
	require.NoError(t, err)
	err = fs.Put(ctx, "test/abc/d/f.txt", NewWriterFromBytes(data1))
	require.NoError(t, err)
	err = fs.Put(ctx, "test/abc/d/g.txt", NewWriterFromBytes(data1))
	require.NoError(t, err)

	// get folder
	reader, err := fs.Get(ctx, "test/abc/d")
	require.NoError(t, err)

	// put folder
	err = fs.Put(ctx, "new", NewWriter(reader.Node))
	require.NoError(t, err)

	err = fs.Delete(ctx, "")
	require.NoError(t, err)
}

func TestStoreIterate(t *testing.T) {
	ctx := context.Background()
	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	// add datas
	data := []byte("test1")
	data2 := []byte("test2")
	data3 := []byte("test3")
	err := fs.Put(ctx, "a/data1.txt", NewWriterFromBytes(data))
	require.NoError(t, err)
	err = fs.Put(ctx, "a/b/c/data2.txt", NewWriterFromBytes(data2))
	require.NoError(t, err)
	err = fs.Put(ctx, "a/b/c/d/e/data3.txt", NewWriterFromBytes(data3))
	require.NoError(t, err)

	// iterate
	err = fs.Iterate("", func(fpath string, reader *Reader) error {
		switch data := reader.Node.(type) {
		case files.Directory:
			fmt.Println(fpath, "is a directory")
		case *files.ReaderFile:
			out, err := io.ReadAll(data)
			require.NoError(t, err)
			fmt.Println(data.AbsPath(), out)
		}
		return nil
	})
	require.NoError(t, err)
}
