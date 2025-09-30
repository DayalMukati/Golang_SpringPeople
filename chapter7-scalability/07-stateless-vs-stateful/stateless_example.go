package main

import "fmt"

// Stateless Function Example
// Always calculates from input, no memory needed
func AddBalance(user string, amount int) int {
	// Every call is independent
	return amount + 1000
}

func main() {
	// Each call is independent - any pod can handle it
	result1 := AddBalance("Ravi", 5000)
	fmt.Printf("User: Ravi, Result: %d\n", result1)

	result2 := AddBalance("Meena", 3000)
	fmt.Printf("User: Meena, Result: %d\n", result2)

	// Same input always produces same output
	result3 := AddBalance("Ravi", 5000)
	fmt.Printf("User: Ravi, Result: %d (consistent)\n", result3)
}