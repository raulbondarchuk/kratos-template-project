// internal/server/middleware/traffic/builder.go
package traffic

import (
	"sync"
	"sync/atomic"

	"golang.org/x/time/rate"
)

// Core limiter components: InFlight and TokenBucket per key

type Builder struct {
	cfg Config

	inflight *inflightLimiter
	rl       rateLimiter
}

// InFlight limiter
type inflightLimiter struct {
	max int64
	cur int64
}

func newInflightLimiter(max int) *inflightLimiter {
	if max < 1 {
		max = 1
	}
	return &inflightLimiter{max: int64(max)}
}

func (l *inflightLimiter) tryEnter() bool {
	for {
		c := atomic.LoadInt64(&l.cur)
		if c >= l.max {
			return false
		}
		if atomic.CompareAndSwapInt64(&l.cur, c, c+1) {
			return true
		}
	}
}

func (l *inflightLimiter) leave() { atomic.AddInt64(&l.cur, -1) }

// Token bucket by key
type tbStore struct {
	mu sync.Mutex
	m  map[string]*rate.Limiter
}

func newTBStore() *tbStore { return &tbStore{m: make(map[string]*rate.Limiter)} }

func (s *tbStore) get(key string, rps float64, burst int) *rate.Limiter {
	if key == "" {
		key = "global"
	}
	if burst < 1 {
		burst = 1
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if lim, ok := s.m[key]; ok {
		return lim
	}
	lim := rate.NewLimiter(rate.Limit(rps), burst)
	s.m[key] = lim
	return lim
}

type rateLimiter interface{ Allow(key string) bool }

type rateLimImpl struct {
	tb *tbStore
	r  float64
	b  int
}

func (r rateLimImpl) Allow(key string) bool { return r.tb.get(key, r.r, r.b).Allow() }

// Builder
func New(cfg Config) *Builder {
	b := &Builder{cfg: cfg}
	if cfg.InflightMax > 0 {
		b.inflight = newInflightLimiter(cfg.InflightMax)
	}
	if cfg.RateRPS > 0 {
		b.rl = rateLimImpl{tb: newTBStore(), r: cfg.RateRPS, b: cfg.RateBurst}
	}
	return b
}

func (b *Builder) tryInflight() (leave func(), blocked bool) {
	if b.inflight == nil {
		return func() {}, false
	}
	if !b.inflight.tryEnter() {
		return func() {}, true
	}
	return func() { b.inflight.leave() }, false
}
