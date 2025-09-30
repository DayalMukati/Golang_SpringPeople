package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// The same event shape the producer sends
type PaymentAuthorized struct {
	Event    string  `json:"event"`
	TxID     string  `json:"tx_id"`
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Merchant string  `json:"merchant_id"`
	When     string  `json:"when"`
}

// A very simple "ledger entry" derived from the event
type LedgerEntry struct {
	TxID      string  `json:"tx_id"`
	Debit     string  `json:"debit"`  // which account is debited
	Credit    string  `json:"credit"` // which account is credited
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	PostedAt  string  `json:"posted_at"`
	Reference string  `json:"reference"` // source event name
}

func main() {
	// Consumer endpoint: POST /events
	http.HandleFunc("/events", handleEvent)

	log.Println("[ledger] listening on :9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}

// handleEvent shows the consumer flow:
// A) Receive the event
// B) Validate the type
// C) Convert to a ledger entry (simple double-entry demo)
// D) "Persist" (for teaching, we just print)
func handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	var ev PaymentAuthorized
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if ev.Event != "payment_authorized" {
		http.Error(w, "unsupported event", http.StatusBadRequest)
		return
	}

	entry := toLedgerEntry(ev)

	// Here you'd insert into a database.
	// For learning, print the JSON so you can see the transformation.
	b, _ := json.MarshalIndent(entry, "", "  ")
	fmt.Println(string(b))

	w.WriteHeader(http.StatusAccepted)
}

func toLedgerEntry(ev PaymentAuthorized) LedgerEntry {
	// NOTE: Example-only accounts to explain the idea:
	// - Debit: customer_clearing (customer's balance reduced)
	// - Credit: merchant_receivable (amount owed to merchant increases)
	return LedgerEntry{
		TxID:      ev.TxID,
		Debit:     "customer_clearing",
		Credit:    "merchant_receivable",
		Amount:    ev.Amount,
		Currency:  ev.Currency,
		PostedAt:  time.Now().UTC().Format(time.RFC3339),
		Reference: ev.Event,
	}
}