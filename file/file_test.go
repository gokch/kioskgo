package file

import (
	"io"
	"testing"

	"github.com/ipfs/boxo/files"
	"github.com/stretchr/testify/require"
)

func TestIPFSFileNew(t *testing.T) {
	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test")
	err := fs.Overwrite("test/abc/d/e.jpg", files.NewBytesFile(data1))
	require.NoError(t, err)

	reader, err := fs.Get("test/abc/d/e.jpg")
	require.NoError(t, err)

	data2, err := io.ReadAll(reader)
	require.Equal(t, data1, data2)
}
