package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read environment variable for feature flag
	featureSplitBill := os.Getenv("FEATURE_SPLIT_BILL") == "true"

	http.HandleFunc("/pay", payHandler(featureSplitBill))
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/readyz", readyHandler)

	log.Println("Payment service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func payHandler(featureSplitBill bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Feature flag: split bill between multiple users
		if featureSplitBill {
			req["feature"] = "split_bill_enabled"
			log.Println("Split bill feature is ENABLED")
		} else {
			req["feature"] = "split_bill_disabled"
			log.Println("Split bill feature is DISABLED")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"request": req,
		})
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// In real scenario, check DB connectivity, external dependencies, etc.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}