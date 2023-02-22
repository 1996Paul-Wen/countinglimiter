package countinglimiter

import (
	"sync"
	"time"
)

// Limiter represents a countinglimiter
type Limiter struct {
	limit        int // the limit of time interval unit
	timeInterval time.Duration

	intervalCumulativeCount int         // the cumulative count in the current time interval
	mu                      *sync.Mutex // protect intervalCumulativeCount from concurrent access
	limited                 bool

	stop         chan struct{} // stop is created when calling Start(), and is closed when calling Stop()
	isChanOpened bool          // mark if stop is opened. true means Limiter is working
}

// NewLimiter generates a Limiter
func NewLimiter(limit int, interval time.Duration) *Limiter {
	return &Limiter{
		limit:        limit,
		timeInterval: interval,
		mu:           &sync.Mutex{},
	}
}

// Start makes Limiter start working
func (l *Limiter) Start() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// already started
	if l.isChanOpened {
		return
	}

	l.stop = make(chan struct{})
	l.isChanOpened = true

	// start a goroutine to monitor stop-signal, and refresh Limiter
	go func() {
		ticker := time.NewTicker(l.timeInterval)
		for {
			select {
			case <-l.stop:
				ticker.Stop()
				return
			case <-ticker.C:
				l.reset()
			}
		}
	}()
}

// reset l.intervalCumulativeCount and l.limited when a new time interval begins
func (l *Limiter) reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.limited = false
	l.intervalCumulativeCount = 0
}

// Allow is the shorthand of AllowN(1)
func (l *Limiter) Allow() bool {
	return l.AllowN(1)
}

// AllowN judges if n requests could pass
func (l *Limiter) AllowN(n int) bool {
	// l is stoped, all can pass
	if !l.isChanOpened {
		return true
	}

	// l is working
	// invalid n, or already limited
	if n <= 0 || l.limited {
		return false
	}

	// note: use read-only scenarios first above to reduce lock contention below

	// lock before calculating
	l.mu.Lock()
	defer l.mu.Unlock()

	temp := l.intervalCumulativeCount + n
	if temp > l.limit {
		return false
	}

	l.intervalCumulativeCount = temp
	l.limited = l.intervalCumulativeCount >= l.limit
	return true
}

// Stop makes Limiter stop and reset the Limiter
func (l *Limiter) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.isChanOpened {
		close(l.stop)
		l.isChanOpened = false
		l.limited = false
		l.intervalCumulativeCount = 0
	}
}
