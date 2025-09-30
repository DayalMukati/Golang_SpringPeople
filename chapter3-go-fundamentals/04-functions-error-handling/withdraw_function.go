package main

import "fmt"

func withdraw(accountNumber string, amount float64) {
	fmt.Println("Withdrawing", amount, "from account", accountNumber)
}

func main() {
	withdraw("ACC12345", 100.0)
}