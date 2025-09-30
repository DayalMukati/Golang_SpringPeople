package main

import (
	"fmt"
	"log"
	"net/http"
)

// A mock Payments service that just echoes the path and request ID.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		fmt.Fprintf(w, `{"service":"payments","path":"%s","request_id":"%s"}`, r.URL.Path, reqID)
	})

	log.Println("[payments] listening on :7002")
	log.Fatal(http.ListenAndServe(":7002", nil))
}