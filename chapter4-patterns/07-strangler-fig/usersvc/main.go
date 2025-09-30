package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// new user service; owns only /api/users/*
func main() {
	http.HandleFunc("/api/users/", handleUsers)

	log.Println("[usersvc] listening on :7001")
	log.Fatal(http.ListenAndServe(":7001", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get("X-Request-ID")

	// Extract user id from path: /api/users/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id string
	if len(parts) >= 3 {
		id = parts[2]
	}

	// Demo response: pretend we read from a NEW datastore
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w,
		`{"service":"usersvc","user_id":"%s","profile":{"name":"Asha","tier":"GOLD"},"request_id":"%s"}`,
		id, reqID)
}