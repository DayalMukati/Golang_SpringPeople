package main

import (
	"fmt"
	"net/http"
)

func checkFraud(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "SAFE") // pretend all transactions are safe
}

func main() {
	http.HandleFunc("/check", checkFraud)
	fmt.Println("Fraud service running on :9090")
	http.ListenAndServe(":9090", nil)
}