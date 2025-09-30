package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type ScoreReq struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type ScoreResp struct {
	Score int    `json:"score"`
	Note  string `json:"note"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/score", func(w http.ResponseWriter, r *http.Request) {
		var in ScoreReq
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		n := rand.Intn(100)
		if n < 20 { // ~20% slow
			time.Sleep(1200 * time.Millisecond)
		} else if n < 35 { // next ~15% 5xx
			http.Error(w, `{"error":"risk_down"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ScoreResp{Score: base(in.Amount), Note: "ok"})
	})

	http.ListenAndServe(":7002", nil)
}

func base(amount float64) int {
	score := 10
	if amount > 1000 {
		score += 40
	}
	if score > 100 {
		score = 100
	}
	return score
}