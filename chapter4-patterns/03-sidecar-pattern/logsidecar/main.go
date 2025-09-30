// logsidecar/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 1) The event we receive from the payments service.
type AuditEvent struct {
	TxID       string  `json:"tx_id"`
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	MerchantID string  `json:"merchant_id"`
	Status     string  `json:"status"`
	Reason     string  `json:"reason,omitempty"`
	IPAddress  string  `json:"ip_address,omitempty"`
	DeviceID   string  `json:"device_id,omitempty"`
}

// 2) The enriched event we will "store".
type EnrichedEvent struct {
	AuditEvent
	Timestamp  string `json:"timestamp"`   // added by sidecar
	FraudScore int    `json:"fraud_score"` // added by sidecar
	Source     string `json:"source"`      // e.g., service name
}

func main() {
	// Register a single endpoint: POST /logs
	http.HandleFunc("/logs", handleLogs)

	log.Println("[sidecar] listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// handleLogs receives one event, enriches it, and prints the final record.
// (In a real app, you would write to a database or Elasticsearch here.)
func handleLogs(w http.ResponseWriter, r *http.Request) {
	// A) Only allow POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	// B) Decode the incoming JSON into our AuditEvent struct.
	var ev AuditEvent
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// C) Enrich the event:
	// - Add UTC timestamp
	// - Compute a simple fraud score
	// - Tag the source (which service sent it)
	enriched := EnrichedEvent{
		AuditEvent: ev,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		FraudScore: simpleFraudScore(ev),
		Source:     "payments-service",
	}

	// D) "Store" the enriched event.
	// For learning, we just print it nicely to the console.
	// (Replace this with real storage later.)
	b, _ := json.MarshalIndent(enriched, "", "  ")
	fmt.Println(string(b))

	// E) Tell the sender we accepted the log.
	w.WriteHeader(http.StatusAccepted)
}

// simpleFraudScore: easy, explainable logic for freshers.
// Start at 0; add points for signals that might be risky.
func simpleFraudScore(ev AuditEvent) int {
	score := 0

	// Higher amounts might carry more risk.
	if ev.Amount > 1000 {
		score += 10
	}

	// If the payment was declined, we add some risk score.
	if ev.Status == "declined" {
		score += 15
	}

	// Non-USD might add a tiny bit (demo purpose only).
	if ev.Currency != "USD" {
		score += 3
	}

	// Missing context can be a small risk.
	if ev.DeviceID == "" || ev.IPAddress == "" {
		score += 2
	}

	return score
}