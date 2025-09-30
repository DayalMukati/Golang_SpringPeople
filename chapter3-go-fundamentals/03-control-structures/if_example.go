package main

import "fmt"

func main() {
	balance := 80.0
	if balance < 100 {
		fmt.Println("Low balance! Please top up your wallet.")
	}
}