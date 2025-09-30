package main

import "fmt"

func main() {
	failedLogins := 2
	if failedLogins > 3 {
		fmt.Println("Account locked due to too many failed logins.")
	} else {
		fmt.Println("Login successful. Welcome!")
	}
}