package main

import (
	"fmt"
	"time"
)

func main() {
	// Step 1: Welcome message
	fmt.Println("Welcome to FinPay Wallet!")

	// Step 2: User name
	userName := "Asha"
	fmt.Println("User:", userName)

	// Step 3: Current date and time
	fmt.Println("Date:", time.Now())

	// Step 4: Wallet balance
	balance := 250.50
	fmt.Println("Wallet Balance: $", balance)
}