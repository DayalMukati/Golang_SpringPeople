package main

import (
	"fmt"
	"time"
)

var cache = make(map[string]int)

func getBalance(user string) int {
	// First check cache
	if val, found := cache[user]; found {
		fmt.Println("Cache hit for", user)
		return val
	}

	// Simulate DB fetch
	fmt.Println("Cache miss for", user, "- fetching from DB...")
	balance := 25000
	cache[user] = balance
	return balance
}

func main() {
	fmt.Println("Balance:", getBalance("Ravi")) // DB fetch
	fmt.Println("Balance:", getBalance("Ravi")) // Cache hit
	time.Sleep(2 * time.Second)
	fmt.Println("Balance:", getBalance("Ravi")) // Cache hit again
}