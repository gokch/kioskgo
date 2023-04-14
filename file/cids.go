package file

import "github.com/ipfs/go-cid"

// cids by relative path
type Cids struct {
	path string
	cids map[cid.Cid]string // map[cid]path
}

func NewCids(path string) *Cids {
	return &Cids{
		path: path,
		cids: map[cid.Cid]string{},
	}
}

func (c *Cids) Get() map[cid.Cid]string {
	return c.cids
}

func (c *Cids) Path() string {
	return c.path
}

func (c *Cids) Add(cid cid.Cid) {
	c.cids[cid] = c.path
}

func (c *Cids) Remove(cid cid.Cid) {
	delete(c.cids, cid)
}

func (c *Cids) Clear() {
	c.cids = map[cid.Cid]string{}
}
