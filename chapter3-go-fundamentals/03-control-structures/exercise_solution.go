package main

import "fmt"

func main() {
	// User wallet balance
	balance := 650.0

	// 1) If-Else check for premium offers
	if balance > 500 {
		fmt.Println("You are eligible for premium offers!")
	} else {
		fmt.Println("Keep using your wallet to unlock premium offers.")
	}

	// 2) For loop for cashback rewards
	fmt.Println("Daily Cashback Rewards:")
	for day := 1; day <= 5; day++ {
		fmt.Println("Day", day, ": $2 credited")
	}
}