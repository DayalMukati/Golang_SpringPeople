package main

import (
	"encoding/json"
	"net/http"
)

type PaymentRequest struct {
	FromAccount string  `json:"fromAccount"`
	ToAccount   string  `json:"toAccount"`
	Amount      float64 `json:"amount"`
}

type PaymentResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transactionId"`
}

func handlePayment(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	// Simulate processing
	res := PaymentResponse{
		Status:        "SUCCESS",
		TransactionID: "TXN12345",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/payments", handlePayment)
	http.ListenAndServe(":8080", nil)
}