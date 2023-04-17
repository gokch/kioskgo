package p2p

import (
	"runtime/debug"
	"sync"

	"github.com/ipfs/go-cid"
)

// 서버일 경우 waitlist 를 우선순위대로 처리
// 클라이언트의 경우 waitlist 를 특정 피어에 요청
// file 쪽으로 가고 지울 놈들임...

type fileManager struct {
	mtx sync.Mutex

	cids  map[cid.Cid]string          // map[cid]localPath
	paths map[string]map[cid.Cid]bool // map[localPath]cids
}

func NewCids() *fileManager {
	return &fileManager{
		cids:  map[cid.Cid]string{},
		paths: map[string]map[cid.Cid]bool{},
	}
}

func (c *fileManager) GetPath(cid cid.Cid) string {
	return c.cids[cid]
}

func (c *fileManager) GetCids(path string) map[cid.Cid]bool {
	return c.paths[path]
}

func (c *fileManager) Add(ci cid.Cid, path string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cids[ci] = path
	if _, ok := c.paths[path]; !ok {
		c.paths[path] = map[cid.Cid]bool{}
	}
	c.paths[path][ci] = true
}

func (c *fileManager) Remove(cid cid.Cid, path string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.cids, cid)
	delete(c.paths, path)
}

func (c *fileManager) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cids = map[cid.Cid]string{}
	c.paths = map[string]map[cid.Cid]bool{}

	// clear orphan memory
	debug.FreeOSMemory()
}
