package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	usersURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	usersProxy := httputil.NewSingleHostReverseProxy(usersURL)

	mux := http.NewServeMux()

	mux.Handle("/api/users/", usersProxy)

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