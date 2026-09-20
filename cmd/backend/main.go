package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"message":"response from users service","path":"%s"}`, r.URL.Path)
	})

	log.Println("Users service listening on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
