# Service Mesh Integration Example

This example demonstrates Service Mesh concepts with retries, mTLS simulation, and tracing.

## Services

1. **payments** - Runs on port 9000
   - Exposes POST /pay
   - Calls mesh proxy (not directly to ledger)
   - No retry/TLS/tracing logic (mesh handles it)

2. **meshproxy** (sidecar-like) - Runs on port 15001
   - Forwards requests to ledger
   - Adds trace IDs (X-Request-ID)
   - Adds "mTLS-like" identity headers
   - Applies timeouts and retries (2 retries on 5xx)
   - Prints metrics

3. **ledger** - Runs on port 7002
   - Exposes POST /ledger/debit
   - Verifies mesh identity headers
   - Randomly delays or errors (to show retries)

## How to Run

1. Start ledger:
```bash
cd ledger
go run main.go
```

2. Start meshproxy:
```bash
cd meshproxy
go run main.go
```

3. Start payments:
```bash
cd payments
go run main.go
```

4. Send a payment:
```bash
curl -X POST http://localhost:9000/pay \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":50,"currency":"USD","merchant_id":"M10"}'
```

5. Observe:
   - Payments responds quickly with ledger result
   - Meshproxy logs trace ID, retries, and status
   - Ledger sometimes delays or returns 5xx (mesh retries automatically)