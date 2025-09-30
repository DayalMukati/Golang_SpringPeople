package main

import "fmt"

func addFee(balance float64) {
	balance -= 5
	fmt.Println("Inside function:", balance)
}

func main() {
	bal := 100.0
	addFee(bal)
	fmt.Println("Outside function:", bal) // still 100.0
}