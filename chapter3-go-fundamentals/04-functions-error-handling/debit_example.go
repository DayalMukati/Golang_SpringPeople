package main

import (
	"errors"
	"fmt"
)

// Debit reduces balance if possible, else returns an error
func Debit(balance, amount float64) (float64, error) {
	if amount > balance {
		return balance, errors.New("insufficient funds")
	}
	return balance - amount, nil
}

func main() {
	balance := 500.0
	amount := 600.0

	newBalance, err := Debit(balance, amount)
	if err != nil {
		fmt.Println("Transaction failed:", err)
	} else {
		fmt.Println("Transaction successful! New Balance:", newBalance)
	}
}