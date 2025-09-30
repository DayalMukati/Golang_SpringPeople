package main

import (
	"fmt"
	"net/http"
)

func payHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Payment of 100 processed successfully!")
}

func main() {
	http.HandleFunc("/pay", payHandler)
	fmt.Println("Payment Service running on port 8080")
	http.ListenAndServe(":8080", nil)
}