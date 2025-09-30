// Order Service (after decomposition)
package main

import (
	"fmt"
	"net/http"
)

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Orders endpoint - Microservice")
}

func main() {
	http.HandleFunc("/orders", ordersHandler)
	fmt.Println("Order service listening on port 8082")
	http.ListenAndServe(":8082", nil)
}