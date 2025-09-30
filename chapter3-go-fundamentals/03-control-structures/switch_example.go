package main

import "fmt"

func main() {
	tier := "Gold"
	switch tier {
	case "Gold":
		fmt.Println("Gold Member: 5% cashback applied.")
	case "Silver":
		fmt.Println("Silver Member: 2% cashback applied.")
	default:
		fmt.Println("Standard Member: Upgrade for more rewards.")
	}
}