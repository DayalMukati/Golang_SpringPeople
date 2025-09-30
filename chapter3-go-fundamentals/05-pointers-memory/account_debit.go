package main

import "fmt"

// Account struct
type Account struct {
	Number  string
	Balance float64
}

// Function that debits money
func debit(acc *Account, amount float64) {
	if acc.Balance < amount {
		fmt.Println("Not enough balance")
		return
	}
	acc.Balance -= amount
	fmt.Println("Debit successful. New balance:", acc.Balance)
}

func main() {
	acc := Account{Number: "ACC12345", Balance: 500}
	debit(&acc, 200) // pass address
	debit(&acc, 400) // second try
}