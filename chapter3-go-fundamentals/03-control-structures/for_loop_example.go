package main

import "fmt"

func main() {
	transactions := []float64{50, -20, -30}
	fmt.Println("Recent Transactions:")
	for i, tx := range transactions {
		fmt.Println("Txn", i+1, ":", tx)
	}
}