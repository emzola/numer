package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	requests      map[string]int
	lock          sync.Mutex
	limit         int
	timeWindow    time.Duration
	cleanupTicker *time.Ticker
}

var limiter = rateLimiter{
	requests:   make(map[string]int),
	limit:      100, // 100 requests per time window
	timeWindow: time.Minute,
}

func init() {
	// Run periodic cleanup of expired rate-limiting records
	limiter.cleanupTicker = time.NewTicker(limiter.timeWindow)
	go func() {
		for range limiter.cleanupTicker.C {
			limiter.cleanup()
		}
	}()
}

func (l *rateLimiter) cleanup() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.requests = make(map[string]int)
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		limiter.lock.Lock()
		defer limiter.lock.Unlock()

		if _, exists := limiter.requests[clientIP]; !exists {
			limiter.requests[clientIP] = 0
		}

		limiter.requests[clientIP]++
		if limiter.requests[clientIP] > limiter.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
