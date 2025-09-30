// User Service (after decomposition)
package main

import (
	"fmt"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Users endpoint - Microservice")
}

func main() {
	http.HandleFunc("/users", usersHandler)
	fmt.Println("User service listening on port 8081")
	http.ListenAndServe(":8081", nil)
}