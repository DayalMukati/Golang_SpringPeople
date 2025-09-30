package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	user := User{ID: "U1001", Name: "Ravi"}
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/user", userHandler)
	http.ListenAndServe(":8081", nil)
}