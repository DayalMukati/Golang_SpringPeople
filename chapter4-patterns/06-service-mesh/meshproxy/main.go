package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"time"
)

/*
meshproxy listens on :15001
- Forwards requests to the real Ledger (http://localhost:7002)
- Adds trace id if missing (X-Request-ID)
- Adds "mTLS-like" identity headers:
    X-Mesh-mTLS: true
    X-Service-Identity: payments
- Applies per-try timeout and retry logic (2 retries on 5xx/timeout)
- Emits tiny metrics to stdout
*/

const (
	upstreamURL   = "http://localhost:7002" // ledger base
	perTryTimeout = 300 * time.Millisecond
	maxRetries    = 2
)

var (
	totalReq   int
	totalRetry int
)

func main() {
	http.HandleFunc("/", handleProxy)

	log.Println("[meshproxy] listening on :15001 (forwarding to", upstreamURL, ")")
	log.Fatal(http.ListenAndServe(":15001", nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	totalReq++

	traceID := r.Header.Get("X-Request-ID")
	if traceID == "" {
		traceID = randomID()
	}

	// Copy body (it may be read multiple times across retries)
	inBody, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()

	// Build upstream URL (keep the same path)
	target := upstreamURL + r.URL.Path

	var lastStatus int
	var lastBody []byte
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			totalRetry++
			log.Printf("[meshproxy] retry %d for %s (trace=%s)", attempt, target, traceID)
		}

		req, _ := http.NewRequest(r.Method, target, bytes.NewReader(inBody))
		req.Header = cloneHeaders(r.Header)

		// Mesh adds identity + tracing headers
		req.Header.Set("X-Request-ID", traceID)
		req.Header.Set("X-Mesh-mTLS", "true")
		req.Header.Set("X-Service-Identity", "payments")

		client := &http.Client{Timeout: perTryTimeout}
		resp, e := client.Do(req)
		if e != nil {
			err = e
			continue // try again
		}

		lastStatus = resp.StatusCode
		lastBody, _ = io.ReadAll(resp.Body)
		resp.Body.Close()

		// Retry on 5xx; otherwise break
		if lastStatus >= 500 {
			continue
		}

		err = nil
		break
	}

	if err != nil {
		http.Error(w, `{"error":"upstream_timeout_or_unreachable"}`, http.StatusBadGateway)
		return
	}

	// Pass back upstream response
	w.Header().Set("X-Request-ID", traceID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(lastStatus)
	w.Write(lastBody)

	// Tiny metrics line
	log.Printf("[meshproxy] req=%d retries=%d status=%d trace=%s", totalReq, totalRetry, lastStatus, traceID)
}

func randomID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func cloneHeaders(h http.Header) http.Header {
	out := http.Header{}
	for k, vv := range h {
		for _, v := range vv {
			out.Add(k, v)
		}
	}
	return out
}