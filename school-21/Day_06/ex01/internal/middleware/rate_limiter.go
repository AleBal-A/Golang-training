package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	mu         sync.Mutex
	visitors   = make(map[string]*visitor)
	rateLimit  = rate.Every(2 * time.Second) // 1 запрос каждые 2 секунды
	burstLimit = 2
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rateLimit, burstLimit)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(10 * time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimiter(next http.Handler) http.Handler {
	go cleanupVisitors()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getVisitor(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
