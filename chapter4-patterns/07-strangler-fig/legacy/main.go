package main

import (
	"fmt"
	"log"
	"net/http"
)

// legacy still handles everything EXCEPT /api/users/*
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w,
			`{"service":"legacy","path":"%s","request_id":"%s"}`,
			r.URL.Path, reqID)
	})

	log.Println("[legacy] listening on :7000")
	log.Fatal(http.ListenAndServe(":7000", nil))
}