package main

import "fmt"

// Simulate sharding: splitting data across multiple "shards"
func main() {
	// Single map (like one DB) - all users in one place
	balances := map[string]int{
		"Alice": 1000,
		"Ravi":  25000,
		"Meena": 500,
	}

	fmt.Println("Single DB (no sharding):")
	fmt.Println(balances)

	// Sharded maps - splitting users across multiple "databases"
	shard1 := map[string]int{"Alice": 1000}
	shard2 := map[string]int{"Ravi": 25000, "Meena": 500}

	fmt.Println("\nSharded DBs:")
	fmt.Println("Shard 1:", shard1)
	fmt.Println("Shard 2:", shard2)

	// Query routing: app must know which shard to hit
	user := "Ravi"
	shard := getShardForUser(user)
	fmt.Printf("\nUser '%s' belongs to Shard %d\n", user, shard)

	if shard == 1 {
		fmt.Printf("Balance: %d (from Shard 1)\n", shard1[user])
	} else {
		fmt.Printf("Balance: %d (from Shard 2)\n", shard2[user])
	}
}

// Simple sharding logic: hash user's first letter
func getShardForUser(user string) int {
	if user[0] <= 'M' {
		return 1
	}
	return 2
}