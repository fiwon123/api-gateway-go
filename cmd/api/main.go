package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

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