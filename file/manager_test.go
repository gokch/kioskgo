package file

import "testing"

func TestManager(t *testing.T) {
	manager := NewFileManager("rootPath")
	_ = manager
}
