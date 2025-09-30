package main

import (
	"fmt"
	"log"
	"net/http"
)

// A mock Users service that just echoes the path and request ID.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		fmt.Fprintf(w, `{"service":"users","path":"%s","request_id":"%s"}`, r.URL.Path, reqID)
	})

	log.Println("[users] listening on :7001")
	log.Fatal(http.ListenAndServe(":7001", nil))
}