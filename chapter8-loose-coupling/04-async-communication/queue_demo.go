package main

import (
	"fmt"
	"time"
)

// Transaction message structure
type Transaction struct {
	ID     string
	Amount float64
}

// Payment Service: sends transactions to the queue
func paymentService(queue chan Transaction) {
	for i := 1; i <= 3; i++ {
		tx := Transaction{ID: fmt.Sprintf("TXN%d", i), Amount: float64(i * 1000)}
		fmt.Println("Payment Service: queued", tx.ID)
		queue <- tx
		time.Sleep(1 * time.Second) // simulate gap between payments
	}
}

// Fraud Service: processes transactions from the queue
func fraudService(queue chan Transaction) {
	for tx := range queue {
		fmt.Println("Fraud Service: processing", tx.ID)
		time.Sleep(2 * time.Second) // simulate slow fraud check
		fmt.Println("Fraud Service: approved", tx.ID)
	}
}

func main() {
	queue := make(chan Transaction, 5) // buffered channel as a queue

	// Start Fraud Service in background
	go fraudService(queue)

	// Run Payment Service
	paymentService(queue)

	time.Sleep(5 * time.Second) // wait for Fraud to finish
}