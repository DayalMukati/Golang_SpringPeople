package main

import "fmt"

type Account struct {
	Number  string
	Balance float64
}

func topUp(acc *Account, amount float64) {
	acc.Balance += amount
}

func main() {
	acc := Account{Number: "ACC200", Balance: 100}
	fmt.Println("Initial balance:", acc.Balance)
	topUp(&acc, 50)
	fmt.Println("After top-up:", acc.Balance)
}