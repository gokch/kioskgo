package file

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileNew(t *testing.T) {
	fs := NewFileSystem("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test")
	err := fs.Add("test/abc/d", "filename2", bytes.NewBuffer(data1))
	require.NoError(t, err)

	file, err := fs.Get("test/abc/d", "filename2")
	require.NoError(t, err)

	data2, err := io.ReadAll(file)
	require.NoError(t, err)

	require.Equal(t, data1, data2)
}
