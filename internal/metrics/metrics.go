package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var Requests atomic.Uint64

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Requests.Add(1)
		next.ServeHTTP(w, r)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	fmt.Fprintf(
		w,
		"# HELP gateway_requests_total Total HTTP requests received by the gateway.\n"+
			"# TYPE gateway_requests_total counter\n"+
			"gateway_requests_total %d\n",
		Requests.Load(),
	)
}
