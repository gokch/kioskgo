package file

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreNewGet(t *testing.T) {
	ctx := context.Background()

	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test")
	err := fs.Put(ctx, "test/abc/d/e.txt", NewWriterFromBytes(data1))
	require.NoError(t, err)

	reader, err := fs.Get(ctx, "test/abc/d/e.txt")
	require.NoError(t, err)

	data2, err := io.ReadAll(reader)
	require.NoError(t, err)

	require.Equal(t, data1, data2)
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
		out, err := io.ReadAll(reader)
		require.NoError(t, err)
		fmt.Println(reader.AbsPath(), out)
		return nil
	})
	require.NoError(t, err)
}
