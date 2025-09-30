// payments/main.go
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// 1) This is the JSON we expect from the client calling /authorize.
type AuthorizeRequest struct {
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	MerchantID string  `json:"merchant_id"`
	IPAddress  string  `json:"ip_address,omitempty"`
	DeviceID   string  `json:"device_id,omitempty"`
}

// 2) This is the simple event we'll send to the sidecar.
// Keep it small: just enough info for auditing.
type AuditEvent struct {
	TxID       string  `json:"tx_id"`
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	MerchantID string  `json:"merchant_id"`
	Status     string  `json:"status"` // "authorized" or "declined"
	Reason     string  `json:"reason,omitempty"`
	IPAddress  string  `json:"ip_address,omitempty"`
	DeviceID   string  `json:"device_id,omitempty"`
}

func main() {
	// Random seed for our tiny simulation
	rand.Seed(time.Now().UnixNano())

	// Register a single endpoint: POST /authorize
	http.HandleFunc("/authorize", authorizeHandler)

	log.Println("[payments] listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// authorizeHandler explains the full flow, step-by-step:
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	// A) Parse the incoming JSON.
	var req AuthorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// B) Make a very simple decision:
	// - If amount <= 1000: authorize.
	// - If amount > 1000: 40% chance we decline (simulating risk).
	status := "authorized"
	reason := ""
	if req.Amount > 1000 && rand.Intn(100) < 40 {
		status = "declined"
		reason = "RISK_RULE"
	}

	// C) Create a transaction ID (here we use time for simplicity).
	txID := time.Now().Format("20060102150405.000")

	// D) Build the audit event we want to send to the sidecar.
	event := AuditEvent{
		TxID:       txID,
		UserID:     req.UserID,
		Amount:     req.Amount,
		Currency:   req.Currency,
		MerchantID: req.MerchantID,
		Status:     status,
		Reason:     reason,
		IPAddress:  req.IPAddress,
		DeviceID:   req.DeviceID,
	}

	// E) Send the event to the sidecar asynchronously.
	// (If the sidecar is down, we just log the error in this simple demo.)
	go sendToSidecar(event)

	// F) Reply to the client with a small JSON.
	resp := map[string]string{
		"tx_id":  txID,
		"status": status,
		"reason": reason,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// sendToSidecar posts the event to http://localhost:9000/logs
func sendToSidecar(ev AuditEvent) {
	body, _ := json.Marshal(ev)
	req, _ := http.NewRequest("POST", "http://localhost:9000/logs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[payments] sidecar not reachable: %v\n", err)
		return
	}
	_ = resp.Body.Close()
}