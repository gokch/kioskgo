package file

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreNewGet(t *testing.T) {
	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test")
	err := fs.Overwrite("test/abc/d/e.txt", NewWriterFromBytes(data1))
	require.NoError(t, err)

	reader, err := fs.Get("test/abc/d/e.txt")
	require.NoError(t, err)

	data2, err := io.ReadAll(reader)
	require.Equal(t, data1, data2)
}

func TestFileIterate(t *testing.T) {
	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	// add datas
	data := []byte("test1")
	data2 := []byte("test2")
	data3 := []byte("test3")
	err := fs.Overwrite("a/data1.txt", NewWriterFromBytes(data))
	require.NoError(t, err)
	err = fs.Overwrite("a/b/c/data2.txt", NewWriterFromBytes(data2))
	require.NoError(t, err)
	err = fs.Overwrite("a/b/c/d/e/data3.txt", NewWriterFromBytes(data3))
	require.NoError(t, err)

	// iterate
	readers, err := fs.Iterate("")
	require.NoError(t, err)
	for _, r := range readers {
		out, err := io.ReadAll(r)
		require.NoError(t, err)
		fmt.Println(r.AbsPath(), out)
	}
}
