package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	http.HandleFunc("/", route)

	log.Println("[proxy] listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// route decides legacy vs new service.
// - /api/users/... -> users service (new)
// - everything else -> legacy
func route(w http.ResponseWriter, r *http.Request) {
	reqID := newReqID()
	w.Header().Set("X-Request-ID", reqID)
	r.Header.Set("X-Request-ID", reqID)

	var target *url.URL
	if strings.HasPrefix(r.URL.Path, "/api/users/") {
		target, _ = url.Parse("http://localhost:7001")
	} else {
		target, _ = url.Parse("http://localhost:7000")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Optional: drop the prefix for the new service if you want
	orig := proxy.Director
	proxy.Director = func(req *http.Request) {
		orig(req)
		req.Header.Set("X-Request-ID", reqID)
		// keep path as-is to keep demo simple
	}

	proxy.ServeHTTP(w, r)
}

func newReqID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}