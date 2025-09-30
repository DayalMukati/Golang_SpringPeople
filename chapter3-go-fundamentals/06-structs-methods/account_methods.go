package main

import "fmt"

// Define Account struct
type Account struct {
	Number  string
	Owner   string
	Balance float64
}

// Method to add money
func (a *Account) TopUp(amount float64) {
	a.Balance += amount
	fmt.Println("Top-up successful. New balance:", a.Balance)
}

// Method to debit money
func (a *Account) Debit(amount float64) {
	if a.Balance < amount {
		fmt.Println("Insufficient funds")
		return
	}
	a.Balance -= amount
	fmt.Println("Debit successful. New balance:", a.Balance)
}

// Method to check balance
func (a Account) CheckBalance() {
	fmt.Println("Current balance:", a.Balance)
}

func main() {
	// Create a new account
	acc := Account{Number: "ACC123", Owner: "Asha", Balance: 500.0}

	// Call methods
	acc.CheckBalance()
	acc.TopUp(200)
	acc.Debit(600)
	acc.CheckBalance()
}