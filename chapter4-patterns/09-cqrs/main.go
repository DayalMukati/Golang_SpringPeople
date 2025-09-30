package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// --- Write model (ledger) ---
type Payment struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"` // "debit" or "credit"
	TxID   string  `json:"tx_id"`
}

// --- Read model (balances) ---
var (
	ledger   []Payment
	balances = map[string]float64{}
	mu       sync.Mutex
)

// Command: POST /command/pay
func handlePay(w http.ResponseWriter, r *http.Request) {
	var p Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Add to ledger (write model)
	ledger = append(ledger, p)

	// Update balance view (read model)
	if p.Type == "credit" {
		balances[p.UserID] += p.Amount
	} else if p.Type == "debit" {
		balances[p.UserID] -= p.Amount
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","tx_id":"%s"}`, p.TxID)
}

// Query: GET /query/balance?user=U1
func handleBalance(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	mu.Lock()
	bal := balances[user]
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"user":"%s","balance":%.2f}`, user, bal)
}

func main() {
	http.HandleFunc("/command/pay", handlePay)
	http.HandleFunc("/query/balance", handleBalance)

	log.Println("[cqrs] listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}