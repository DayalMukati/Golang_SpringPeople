package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

/*
payments exposes POST /pay
- Parses a simple payment request
- Calls the local mesh proxy (http://localhost:15001/ledger/debit)
- Does not implement retries/TLS/tracing (the mesh does that)
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
	http.HandleFunc("/pay", handlePay)

	log.Println("[payments] listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handlePay(w http.ResponseWriter, r *http.Request) {
	var req PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Build a request to the MESH (not directly to ledger)
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "http://localhost:15001/ledger/debit", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	// Small client timeout; mesh will do its own per-try timeout & retries.
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, `{"error":"mesh_unreachable"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Relay response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	_, _ = w.Write(buf.Bytes())
}