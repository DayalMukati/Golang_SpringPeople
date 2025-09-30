package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func payHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:9090/check")
	if err != nil {
		http.Error(w, "Fraud service unavailable", http.StatusServiceUnavailable)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "Payment processed. Fraud check result: %s\n", body)
}

func main() {
	http.HandleFunc("/pay", payHandler)
	fmt.Println("Payment service running on :8080")
	http.ListenAndServe(":8080", nil)
}