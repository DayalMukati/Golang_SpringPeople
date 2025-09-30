package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

const riskURL = "http://localhost:7002/score"

type PayReq struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type Decision struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
	Mode     string `json:"mode"` // "normal" or "fallback"
}

/* Tiny Circuit Breaker (aligned):
- Trip OPEN after 3 consecutive failures
- Stay OPEN 10s (fast-fail)
- After 10s allow 1 probe (HALF-OPEN): success -> CLOSED, fail -> OPEN
*/
type breaker struct {
	mu        sync.Mutex
	state     string // "closed","open","half"
	failCount int
	openUntil time.Time
}

func (b *breaker) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if b.state == "open" {
		if now.After(b.openUntil) {
			b.state = "half"
			return true
		}
		return false
	}
	return true // closed or half
}

func (b *breaker) report(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == "half" {
		if err == nil {
			b.state, b.failCount = "closed", 0
		} else {
			b.state, b.openUntil = "open", time.Now().Add(10*time.Second)
		}
		return
	}

	if err == nil {
		b.failCount = 0
		return
	}

	b.failCount++
	if b.failCount >= 3 {
		b.state, b.openUntil = "open", time.Now().Add(10*time.Second)
	}
}

var brk = &breaker{state: "closed"}

func main() {
	http.HandleFunc("/pay", handlePay)
	http.ListenAndServe(":9000", nil)
}

func handlePay(w http.ResponseWriter, r *http.Request) {
	var in PayReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	// 1) Ask breaker if we should try Risk
	if !brk.allow() {
		json.NewEncoder(w).Encode(fallback(in))
		return
	}

	// 2) Call Risk with a short timeout (600ms)
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest("POST", riskURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 600 * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		brk.report(err)
		json.NewEncoder(w).Encode(fallback(in))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		brk.report(errors.New("5xx"))
		json.NewEncoder(w).Encode(fallback(in))
		return
	}

	// Success path
	brk.report(nil)
	json.NewEncoder(w).Encode(Decision{Approved: true, Reason: "risk_ok", Mode: "normal"})
}

func fallback(in PayReq) Decision {
	// Simple, safe:
	if in.Amount <= 50 {
		return Decision{Approved: true, Reason: "challenge_small", Mode: "fallback"}
	}
	return Decision{Approved: false, Reason: "hold_high_amount", Mode: "fallback"}
}