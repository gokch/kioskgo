package file

import "testing"

func TestMFS(t *testing.T) {
	mfs := NewMfs(NewFileStore("rootpath"))
	rootDir := mfs.Root.GetDirectory()
	rootDir.Mkdir("tesbeset")
}
