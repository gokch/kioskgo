package p2p

import (
	"sync"

	"github.com/ipfs/go-cid"
)

// 서버일 경우 waitlist 를 우선순위대로 처리
// 클라이언트의 경우 waitlist 를 특정 피어에 요청

type cids struct {
	cids sync.Map
}

func (w *cids) Add(cid cid.Cid) {

}

func (w *cids) Remove() {

}

func (w *cids) Clear() {

}
