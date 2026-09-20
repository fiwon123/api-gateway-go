package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	usersURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	expectedToken := os.Getenv("API_TOKEN")
	if expectedToken == "" {
		log.Fatal("API_TOKEN is required")
	}

	usersProxy := httputil.NewSingleHostReverseProxy(usersURL)

	mux := http.NewServeMux()

	mux.Handle("/api/users/", withMiddleware(
							usersProxy,
							requestID,
							logging,
							authentication(expectedToken),
						))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Println("API gateway listening on http://localhost:8080")
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}


}

type middleware func(http.Handler) http.Handler

func withMiddleware(handler http.Handler, middlewares ...middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")

		if id == "" {
			id = newRequestID()
		}

		r.Header.Set("X-Request-ID", id)
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"method=%s path=%s duration=%s request_id=%s",
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.Header.Get("X-Request-ID"),
		)
	})
}

func newRequestID() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(bytes)
}

func authentication(expectedToken string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader,  "Bearer ") {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(
					w,
					`{"error":"missing bearer token"}`,
					http.StatusUnauthorized,
				)
				return
			}

			providedToken := strings.TrimPrefix(authHeader, "Bearer ")

			if subtle.ConstantTimeCompare(
				[]byte(providedToken),
				[]byte(expectedToken),
			) != 1 {
				http.Error(
					w,
					`{"error":"invalid token"}`,
					http.StatusUnauthorized,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}