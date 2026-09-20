package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiwon123/api-gateway-go/internal/auth"
	"github.com/fiwon123/api-gateway-go/internal/limiter"
	"github.com/fiwon123/api-gateway-go/internal/middleware"
)

func main() {
	usersURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	orderURL, err := url.Parse("http://localhost:8082")
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	usersProxy := newProxy(usersURL)
	ordersProxy := newProxy(orderURL)

	mux := http.NewServeMux()

	mux.Handle("/api/users/", middleware.WithMiddleware(
		usersProxy,
		requestID,
		logging,
		limiter.RateLimit(60, time.Minute),
		auth.JwtAuthentication(
			[]byte(jwtSecret),
			os.Getenv("JWT_ISSUER"),
			os.Getenv("JWT_AUDIENCE"),
		),
	))

	mux.Handle(
		"/api/orders/",
		middleware.WithMiddleware(
			ordersProxy,
			requestID,
			logging,
			limiter.RateLimit(60, time.Minute),
			auth.JwtAuthentication([]byte(jwtSecret),
				os.Getenv("JWT_ISSUER"),
				os.Getenv("JWT_AUDIENCE")),
		),
	)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("API gateway listening on http://localhost:8080")

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}

	case signal := <-shutdown:
		log.Printf("received signal: %s", signal)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	log.Println("shutting down gateway")

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("gateway stopped")
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

func newProxy(target *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf(
			"backend error path=%s request_id=%s error=%v",
			r.URL.Path,
			r.Header.Get("X-Request-ID"),
			err,
		)

		w.Header().Set("Content-Type", "application/json")
		http.Error(
			w,
			`{"error":"backend unavailable"}`,
			http.StatusBadGateway,
		)
	}

	return proxy
}
