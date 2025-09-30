package main

import (
	"encoding/json"
	"net/http"
)

type PaymentResponse struct {
	Status string `json:"status"`
	Amount int    `json:"amount"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	resp := PaymentResponse{Status: "Payment Successful", Amount: 100}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/pay", paymentHandler)
	http.ListenAndServe(":8082", nil)
}