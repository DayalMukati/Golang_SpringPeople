package main

import "fmt"

func withdraw(balance, amount float64) float64 {
	newBalance := balance - amount
	return newBalance
}

func main() {
	balance := 500.0
	newBalance := withdraw(balance, 100.0)
	fmt.Println("New balance:", newBalance)
}