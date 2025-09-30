package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func checkFraud(w http.ResponseWriter, r *http.Request) {
	// Simulate intermittent failures (30% failure rate)
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < 30 {
		fmt.Println("Simulating fraud service failure...")
		time.Sleep(3 * time.Second) // Simulate timeout
		return
	}

	fmt.Fprintln(w, "SAFE")
}

func main() {
	http.HandleFunc("/check", checkFraud)
	fmt.Println("Fraud service (with simulated failures) running on :9090")
	http.ListenAndServe(":9090", nil)
}