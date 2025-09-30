package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func payHandler(w http.ResponseWriter, r *http.Request) {
	// First attempt to call Fraud Service
	resp, err := http.Get("http://localhost:9090/check")

	if err != nil {
		// Retry once before failing
		fmt.Println("First attempt failed, retrying...")
		time.Sleep(500 * time.Millisecond)

		resp, err = http.Get("http://localhost:9090/check")
		if err != nil {
			http.Error(w, "Fraud service unavailable after retry", http.StatusServiceUnavailable)
			return
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "Payment processed. Fraud check result: %s\n", body)
}

func main() {
	http.HandleFunc("/pay", payHandler)
	fmt.Println("Payment service with retry running on :8080")
	http.ListenAndServe(":8080", nil)
}