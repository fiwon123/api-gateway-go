package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/orders/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")

		fmt.Fprintf(
			w,
			`{"message":"response from orders service","path":"%s","user_id":"%s"}`,
			r.URL.Path,
			userID,
		)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"orders"}`))
	})

	log.Println("Orders service listening on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
