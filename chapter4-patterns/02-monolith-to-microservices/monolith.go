// Monolithic server (before decomposition)
package main

import (
	"fmt"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Users endpoint - Monolith")
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Orders endpoint - Monolith")
}

func main() {
	// Both users and orders handled in one application
	http.HandleFunc("/users", usersHandler)  // e.g., list or get users
	http.HandleFunc("/orders", ordersHandler) // e.g., create or get orders

	fmt.Println("Monolith listening on port 8080")
	http.ListenAndServe(":8080", nil)
}