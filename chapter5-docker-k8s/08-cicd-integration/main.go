package main

import (
	"fmt"
	"net/http"
)

func payHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Payment processed successfully!")
}

func main() {
	http.HandleFunc("/pay", payHandler)
	fmt.Println("Payment service running on :8080")
	http.ListenAndServe(":8080", nil)
}