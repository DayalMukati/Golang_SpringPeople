package main

import (
	"errors"
	"fmt"
)

func withdraw(balance, amount float64) (float64, error) {
	if balance < amount {
		return balance, errors.New("insufficient funds")
	}
	return balance - amount, nil
}

func main() {
	newBalance, err := withdraw(100, 150)
	if err != nil {
		fmt.Println("Payment failed:", err)
	} else {
		fmt.Println("Payment success! New balance:", newBalance)
	}
}