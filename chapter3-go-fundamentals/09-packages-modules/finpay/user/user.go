package user

import "fmt"

// User struct represents a customer
type User struct {
	ID   string
	Name string
}

// Method to greet the user
func (u User) Greet() {
	fmt.Println("Welcome,", u.Name)
}