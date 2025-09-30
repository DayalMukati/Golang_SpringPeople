package main

import (
	"fmt"
	"time"
)

func main() {
	// User details
	userName := "Asha"
	var balance float64 = 250.50
	var rewardPoints int = 120
	var kycDone bool = true
	createdAt := time.Now()

	// Exercise additions
	var accountNumber string = "ACC12345"
	var failedLogins int = 2

	// Print values
	fmt.Println("Welcome to FinPay Wallet!")
	fmt.Println("User:", userName)
	fmt.Println("Account Number:", accountNumber)
	fmt.Println("Wallet Balance: $", balance)
	fmt.Println("Reward Points:", rewardPoints)
	fmt.Println("KYC Completed:", kycDone)
	fmt.Println("Account Created At:", createdAt)
	fmt.Println("Failed Login Attempts:", failedLogins)
}