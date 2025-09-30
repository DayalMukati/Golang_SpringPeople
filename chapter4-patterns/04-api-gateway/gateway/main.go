package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

// A simple API Gateway that:
// 1) Validates an API key
// 2) Applies rate limits
// 3) Adds a request ID
// 4) Routes requests to backend services
// 5) Proxies the request and returns the response

const validAPIKey = "demo-key-123"

// Route table: path prefix -> backend service
var backends = map[string]string{
	"/users/":    "http://localhost:7001",
	"/payments/": "http://localhost:7002",
}

// --- Rate limiting (per IP) ---
type bucket struct {
	tokens int
	last   time.Time
}

var (
	mu      sync.Mutex
	buckets = map[string]*bucket{}
	rate    = 5
	per     = 10 * time.Second
)

func allow(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	b, ok := buckets[ip]
	if !ok {
		b = &bucket{tokens: rate, last: now}
		buckets[ip] = b
	}

	// Reset tokens if enough time passed
	if now.Sub(b.last) >= per {
		b.tokens = rate
		b.last = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// --- Helpers ---
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func randomRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func chooseBackend(path string) *url.URL {
	for prefix, raw := range backends {
		if strings.HasPrefix(path, prefix) {
			u, _ := url.Parse(raw)
			return u
		}
	}
	return nil
}

func main() {
	http.HandleFunc("/", handleGateway)
	log.Println("[gateway] listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleGateway(w http.ResponseWriter, r *http.Request) {
	// 1) API key check
	if r.Header.Get("X-API-Key") != validAPIKey {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// 2) Rate limit
	ip := getClientIP(r)
	if !allow(ip) {
		http.Error(w, `{"error":"rate_limit_exceeded"}`, http.StatusTooManyRequests)
		return
	}

	// 3) Add a request ID
	reqID := randomRequestID()
	w.Header().Set("X-Request-ID", reqID)
	r.Header.Set("X-Request-ID", reqID)

	// 4) Find the backend
	target := chooseBackend(r.URL.Path)
	if target == nil {
		http.Error(w, `{"error":"not_found"}`, http.StatusNotFound)
		return
	}

	// 5) Proxy the request
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Rewrite path: remove prefix before forwarding
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.Header.Set("X-Request-ID", reqID)

		for prefix := range backends {
			if strings.HasPrefix(req.URL.Path, prefix) {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
				if req.URL.Path == "" {
					req.URL.Path = "/"
				}
				break
			}
		}
	}

	proxy.ServeHTTP(w, r)
}