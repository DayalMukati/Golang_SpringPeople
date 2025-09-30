package main

import "fmt"

// Loan struct
type Loan struct {
	LoanID   string
	Borrower string
	Amount   float64
}

// Method to disburse loan
func (l *Loan) Disburse(amount float64) {
	l.Amount += amount
	fmt.Println("Loan disbursed. Total loan amount:", l.Amount)
}

// Method to repay loan
func (l *Loan) Repay(amount float64) {
	if amount > l.Amount {
		fmt.Println("Repayment exceeds loan amount")
		return
	}
	l.Amount -= amount
	fmt.Println("Repayment successful. Remaining loan:", l.Amount)
}

// Method to check status
func (l Loan) CheckStatus() {
	fmt.Println("Loan ID:", l.LoanID)
	fmt.Println("Borrower:", l.Borrower)
	fmt.Println("Remaining Loan Amount:", l.Amount)
}

func main() {
	loan := Loan{LoanID: "LN001", Borrower: "Ravi", Amount: 0}

	loan.CheckStatus()
	loan.Disburse(10000)
	loan.Repay(3000)
	loan.CheckStatus()
}