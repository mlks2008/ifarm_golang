package gpool

import (
	"log"
	"sync"

	"github.com/panjf2000/ants/v2"
)

const (
	PoolMaxSize = 50_0000
)

var (
	PoolIns  *ants.Pool
	syncOnce sync.Once
)

func init() {
	SetGoPoolInstance()
}

// 实例化之后，记得一定要在主协程中回收 PoolIns.Release()
func SetGoPoolInstance() *ants.Pool {
	if PoolIns == nil {
		syncOnce.Do(func() {
			PoolIns, _ = ants.NewPool(PoolMaxSize)
		})
	}
	return PoolIns
}

func Release() {
	if PoolIns == nil {
		return
	}
	PoolIns.Release()
	log.Printf("[SafeGo] pool release")
}

func SafeGo(fn func()) {
	if err := PoolIns.Submit(fn); err != nil {
		log.Printf("[SafeGo] err: %v", err)
	}
}
