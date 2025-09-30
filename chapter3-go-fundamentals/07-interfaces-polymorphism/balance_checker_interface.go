package main

import "fmt"

// Interface
type BalanceChecker interface {
	CheckBalance()
}

// Account struct
type Account struct {
	Number  string
	Balance float64
}

func (a Account) CheckBalance() {
	fmt.Println("Account", a.Number, "has balance:", a.Balance)
}

// Loan struct
type Loan struct {
	LoanID  string
	Balance float64
}

func (l Loan) CheckBalance() {
	fmt.Println("Loan", l.LoanID, "outstanding amount:", l.Balance)
}

// Function that works on any BalanceChecker
func PrintBalance(b BalanceChecker) {
	b.CheckBalance()
}

func main() {
	acc := Account{Number: "ACC100", Balance: 300}
	loan := Loan{LoanID: "LN200", Balance: 5000}

	// Both implement BalanceChecker
	PrintBalance(acc)
	PrintBalance(loan)
}