package main

import "fmt"

func deductFee(balance *float64) {
	*balance -= 10 // update the original value
}

func main() {
	myBalance := 100.0
	deductFee(&myBalance) // pass the address
	fmt.Println("After fee:", myBalance) // Now 90
}