package main

import "fmt"

func main() {
	userName := "Asha"
	balance := 80.0
	failedLogins := 2
	tier := "Gold"

	fmt.Println("Welcome,", userName)

	// If: check balance
	if balance < 100 {
		fmt.Println("⚠️ Low balance! Please top up your wallet.")
	}

	// If-Else: check login attempts
	if failedLogins > 3 {
		fmt.Println("Account locked due to too many failed logins.")
	} else {
		fmt.Println("Login successful.")
	}

	// Switch: membership tier
	switch tier {
	case "Gold":
		fmt.Println("Gold Member: 5% cashback applied.")
	case "Silver":
		fmt.Println("Silver Member: 2% cashback applied.")
	default:
		fmt.Println("Standard Member: Upgrade for more rewards.")
	}

	// For loop: transaction history
	transactions := []float64{50, -20, -30}
	fmt.Println("Recent Transactions:")
	for i, tx := range transactions {
		fmt.Println("Txn", i+1, ":", tx)
	}
}