package main

import (
	"fmt"
	"os"
)

func main() {
	// Simulate discovery: Fraud service address comes from config/env
	fraudService := os.Getenv("FRAUD_SERVICE_URL")

	if fraudService == "" {
		fraudService = "http://fraud-service:8080" // default
	}

	fmt.Println("Payment Service discovered Fraud at:", fraudService)
	// Normally here we would call fraudService via HTTP
}