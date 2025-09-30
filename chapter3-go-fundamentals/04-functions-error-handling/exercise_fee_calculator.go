package main

import (
	"errors"
	"fmt"
)

func calcFee(amount float64) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be > 0")
	}
	return amount * 0.01, nil
}

func main() {
	fee, err := calcFee(100.0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Fee:", fee)
	}

	fee2, err2 := calcFee(-50.0)
	if err2 != nil {
		fmt.Println("Error:", err2)
	} else {
		fmt.Println("Fee:", fee2)
	}
}