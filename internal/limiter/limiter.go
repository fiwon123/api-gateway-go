package limiter

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fiwon123/api-gateway-go/internal/middleware"
)

type visitor struct {
	windowStart time.Time
	requests    int
}

func RateLimit(maxRequests int, window time.Duration) middleware.Middleware {
	var mu sync.Mutex
	visitors := make(map[string]visitor)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client := clientKey(r)
			now := time.Now()

			mu.Lock()

			current, exists := visitors[client]

			if !exists || now.Sub(current.windowStart) >= window {
				current = visitor{
					windowStart: now,
					requests:    0,
				}
			}

			current.requests++
			visitors[client] = current

			allowed := current.requests <= maxRequests

			mu.Unlock()

			if !allowed {
				retryAfter := int(window.Seconds())

				w.Header().Set(
					"Retry-After",
					strconv.Itoa(retryAfter),
				)

				http.Error(
					w,
					`{"error":"rate limit exceeded"}`,
					http.StatusTooManyRequests,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
