package main

import "fmt"

type Account struct {
	Number  string
	Balance float64
}

// Fixed version - uses pointer receiver
func freezeAccount(acc *Account) {
	acc.Balance = 0
}

func main() {
	myAcc := Account{Number: "ACC100", Balance: 200}
	freezeAccount(&myAcc)
	fmt.Println("Balance after freeze:", myAcc.Balance) // now 0
}