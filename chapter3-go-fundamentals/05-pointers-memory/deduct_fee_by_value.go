package main

import "fmt"

func deductFee(balance float64) {
	balance -= 10 // works only on the copy
}

func main() {
	myBalance := 100.0
	deductFee(myBalance)
	fmt.Println("After fee:", myBalance) // Still 100
}