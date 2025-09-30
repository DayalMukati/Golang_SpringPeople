package main

import (
	"fmt"
	"time"
)

// Simulate debiting money
func debit(amount float64, done chan string) {
	time.Sleep(1 * time.Second) // pretend to work
	done <- fmt.Sprintf("Debited: %.2f", amount)
}

// Simulate logging
func logTransaction(amount float64, done chan string) {
	time.Sleep(2 * time.Second)
	done <- fmt.Sprintf("Transaction logged: %.2f", amount)
}

// Simulate sending notification
func notifyUser(amount float64, done chan string) {
	time.Sleep(1 * time.Second)
	done <- fmt.Sprintf("Notification sent for: %.2f", amount)
}

func main() {
	amount := 100.0
	done := make(chan string, 3) // channel to collect results

	// Run all tasks concurrently
	go debit(amount, done)
	go logTransaction(amount, done)
	go notifyUser(amount, done)

	// Collect results
	for i := 0; i < 3; i++ {
		fmt.Println(<-done)
	}
}