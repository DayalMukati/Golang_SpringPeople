package main

import (
	"fmt"
	"net/http"
	"time"
)

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate slow call to Fraud Service
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w, "Payment processed successfully")
}

func main() {
	http.HandleFunc("/pay", paymentHandler)
	http.ListenAndServe(":8080", nil)
}