package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Client sends this JSON to /authorize
type AuthorizeRequest struct {
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	MerchantID string  `json:"merchant_id"`
}

// This is the event we emit when a payment is authorized
type PaymentAuthorized struct {
	Event    string  `json:"event"` // always "payment_authorized"
	TxID     string  `json:"tx_id"`
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Merchant string  `json:"merchant_id"`
	When     string  `json:"when"` // ISO timestamp
}

const ledgerURL = "http://localhost:9001/events" // our consumer's endpoint

func main() {
	rand.Seed(time.Now().UnixNano())

	// One endpoint: POST /authorize
	http.HandleFunc("/authorize", authorizeHandler)

	log.Println("[payments] listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// authorizeHandler:
// A) Parse input
// B) Make a tiny decision (approve most, decline some)
// C) If approved, emit PaymentAuthorized event (async) and respond quickly
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Simple decision rule for learning:
	// - Approve amounts <= 1000
	// - For >1000, approve ~70% of the time (simulate risk)
	approved := req.Amount <= 1000 || rand.Intn(100) < 70

	// Generate a simple transaction ID from the clock (readable for demos)
	txID := time.Now().Format("20060102150405.000")

	// If approved, build and emit the event asynchronously (non-blocking)
	if approved {
		ev := PaymentAuthorized{
			Event:    "payment_authorized",
			TxID:     txID,
			UserID:   req.UserID,
			Amount:   req.Amount,
			Currency: req.Currency,
			Merchant: req.MerchantID,
			When:     time.Now().UTC().Format(time.RFC3339),
		}
		go publishEvent(ev) // <-- decouples producer from consumer speed
	}

	// Respond to the client quickly (producer doesn't wait for consumers)
	resp := map[string]any{
		"tx_id":    txID,
		"status":   ternary(approved, "authorized", "declined"),
		"approved": approved,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// publishEvent sends the event to the ledger's /events endpoint
func publishEvent(ev PaymentAuthorized) {
	body, _ := json.Marshal(ev)
	req, _ := http.NewRequest("POST", ledgerURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Short timeout; in real life a broker gives buffering & retries
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[payments] failed to emit event: %v", err)
		return
	}
	_ = resp.Body.Close()
}

// tiny helper to keep response code neat
func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}