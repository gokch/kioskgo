package p2p

import (
	"runtime/debug"
	"sync"

	"github.com/ipfs/go-cid"
)

// 서버일 경우 waitlist 를 우선순위대로 처리
// 클라이언트의 경우 waitlist 를 특정 피어에 요청

type Cids struct {
	mtx sync.Mutex

	cids  map[cid.Cid]string // map[cid]localPath
	paths map[string]cid.Cid // map[localPath]cid
}

func NewCids() *Cids {
	return &Cids{
		cids:  map[cid.Cid]string{},
		paths: map[string]cid.Cid{},
	}
}

func (c *Cids) GetPath(cid cid.Cid) string {
	return c.cids[cid]
}

func (c *Cids) GetCid(path string) cid.Cid {
	return c.paths[path]
}

func (c *Cids) Add(cid cid.Cid, path string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cids[cid] = path
	c.paths[path] = cid
}

func (c *Cids) Remove(cid cid.Cid, path string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.cids, cid)
	delete(c.paths, path)
}

func (c *Cids) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cids = map[cid.Cid]string{}
	c.paths = map[string]cid.Cid{}

	// clear orphan memory
	debug.FreeOSMemory()
}
