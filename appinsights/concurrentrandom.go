package appinsights

import (
	"encoding/base64"
	"math/rand"
	"sync/atomic"
	"time"
	"unsafe"
)

type concurrentRandom chan string
var randomGenerator *concurrentRandom

func newConcurrentRandom() *concurrentRandom {
	result := make(concurrentRandom, 4)
	return &result
}

func (generator concurrentRandom) run() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	buf := make([]byte, 9)
	for {
		for i := 0; i < 16384; i++ {
			random.Read(buf)
			generator <- base64.StdEncoding.EncodeToString(buf)
		}
		
		// Reseed every few thousand requests
		random.Seed(time.Now().UnixNano())
	}
}

func randomId() string {
	if randomGenerator == nil {
		r := newConcurrentRandom()
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&randomGenerator)), unsafe.Pointer(nil), unsafe.Pointer(r)) {
			go r.run()
		} else {
			close(*r)
		}
	}
	
	return <- *randomGenerator
}
