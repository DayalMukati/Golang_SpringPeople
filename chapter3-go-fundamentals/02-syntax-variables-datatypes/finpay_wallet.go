package main

import (
	"fmt"
	"time"
)

func main() {
	// String variable
	userName := "Asha"

	// Float variable (wallet balance)
	var balance float64 = 250.50

	// Integer variable (reward points)
	var rewardPoints int = 120

	// Boolean variable (KYC completed or not)
	var kycDone bool = true

	// Time variable (account created date)
	createdAt := time.Now()

	// Print values
	fmt.Println("Welcome to FinPay Wallet!")
	fmt.Println("User:", userName)
	fmt.Println("Wallet Balance: $", balance)
	fmt.Println("Reward Points:", rewardPoints)
	fmt.Println("KYC Completed:", kycDone)
	fmt.Println("Account Created At:", createdAt)
}