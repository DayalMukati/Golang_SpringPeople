package main

import "fmt"

// Stateful Example (not ideal for scaling)
var balance = 1000

func UpdateBalance(amount int) int {
	// Modifies global state - remembers previous calls
	balance += amount
	return balance
}

func main() {
	// Each call modifies shared state
	fmt.Printf("Initial balance: %d\n", balance)

	result1 := UpdateBalance(500)
	fmt.Printf("After adding 500: %d\n", result1)

	result2 := UpdateBalance(300)
	fmt.Printf("After adding 300: %d\n", result2)

	// Problem: If this runs in multiple pods,
	// keeping balances consistent becomes a headache
	fmt.Println("\nNote: This pattern doesn't scale well across multiple pods!")
}