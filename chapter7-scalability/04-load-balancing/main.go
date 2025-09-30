package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Reply struct {
	Pod       string `json:"pod"`
	Timestamp string `json:"ts"`
	Path      string `json:"path"`
}

func pay(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	resp := Reply{
		Pod:       host,
		Timestamp: time.Now().Format(time.RFC3339),
		Path:      r.URL.Path,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/pay", pay)
	http.ListenAndServe(":8080", nil)
}