package file

import (
	"testing"
)

func TestManager(t *testing.T) {
	manager := NewFileManager("rootPath")
	_ = manager
	// ci := manager.extractCid(".//a/b//c/0xaaaaa")
	// fmt.Println("cid :", ci)
}
