package main

import (
	"finpay/payment"
	"finpay/user"
)

func main() {
	// Create a new user
	u := user.User{ID: "U1001", Name: "Ravi"}
	u.Greet()

	// Process a payment
	payment.Process(250.0)
}