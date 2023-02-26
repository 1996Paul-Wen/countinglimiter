package countinglimiter

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	l := NewLimiter(1000, 1*time.Second)

	// stoped limiter case: 0 head off
	var headOff int64 = 0
	for i := 0; i < 2000; i++ {
		if !l.AllowN(1) {
			headOff += 1
		}
	}
	fmt.Printf("stoped limiter case. headOff: %d\n", headOff) // 0

	// started limiter case
	l.Start()
	headOff = 0
	wg := &sync.WaitGroup{}
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Allow() {
				atomic.AddInt64(&headOff, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("started limiter case. headOff: %d\n", headOff) // 1000

	// stop and restart case
	l.Stop()
	l.Start()
	headOff = 0
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Allow() {
				atomic.AddInt64(&headOff, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("stop and restart case. headOff: %d\n", headOff) // 1000
}
