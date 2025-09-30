package main

import (
	"errors"
	"fmt"
)

var ErrManualReview = errors.New("requires manual review")

func validateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be > 0")
	}
	return nil
}

func checkAndDebit(balance, amount float64) (float64, error) {
	if balance < amount {
		return balance, errors.New("insufficient funds")
	}
	return balance - amount, nil
}

func authorizeAndDebit(balance, amount float64) (float64, error) {
	if err := validateAmount(amount); err != nil {
		return balance, err
	}
	if amount > 1000 {
		return balance, ErrManualReview
	}
	return checkAndDebit(balance, amount)
}

func main() {
	// Test case 1: Normal transaction
	newBalance, err := authorizeAndDebit(500, 200)
	if err != nil {
		fmt.Println("Transaction failed:", err)
	} else {
		fmt.Println("Transaction successful! New balance:", newBalance)
	}

	// Test case 2: Large amount requiring review
	newBalance2, err2 := authorizeAndDebit(2000, 1500)
	if err2 != nil {
		fmt.Println("Transaction failed:", err2)
	} else {
		fmt.Println("Transaction successful! New balance:", newBalance2)
	}
}