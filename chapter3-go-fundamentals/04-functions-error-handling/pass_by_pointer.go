package main

import "fmt"

func addFee(balance *float64) {
	*balance -= 5
}

func main() {
	bal := 100.0
	addFee(&bal) // pass address of variable
	fmt.Println("New balance:", bal) // now 95
}