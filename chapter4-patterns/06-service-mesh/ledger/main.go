package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/*
ledger exposes POST /ledger/debit
- Verifies "mTLS-like" headers (simulating identity from the mesh)
- Randomly delays or errors to demonstrate mesh retries
- Returns a simple JSON result on success
*/

type PayRequest struct {
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Merchant string  `json:"merchant_id"`
}

type LedgerResult struct {
	Status   string  `json:"status"`
	TxID     string  `json:"tx_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/ledger/debit", handleDebit)

	log.Println("[ledger] listening on :7002")
	log.Fatal(http.ListenAndServe(":7002", nil))
}

func handleDebit(w http.ResponseWriter, r *http.Request) {
	// Simulate mTLS verification via headers (demo only)
	if r.Header.Get("X-Mesh-mTLS") != "true" ||
		r.Header.Get("X-Service-Identity") != "payments" {
		http.Error(w, `{"error":"unauthenticated_mesh"}`, http.StatusUnauthorized)
		return
	}

	// Random slowness / failure to show retries
	n := rand.Intn(100)
	switch {
	case n < 20:
		time.Sleep(500 * time.Millisecond) // slower than mesh per-try timeout (forces retry)
	case n >= 20 && n < 30:
		http.Error(w, `{"error":"temporary_storage_error"}`, http.StatusInternalServerError)
		return
	}

	var in PayRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	out := LedgerResult{
		Status:   "posted",
		TxID:     time.Now().Format("20060102150405.000"),
		Amount:   in.Amount,
		Currency: in.Currency,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}