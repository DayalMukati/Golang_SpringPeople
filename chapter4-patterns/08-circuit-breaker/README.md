# Circuit Breaker & Fault Tolerance Example

This example demonstrates the Circuit Breaker Pattern with fallback handling.

## Services

1. **risk** (downstream) - Runs on port 7002
   - Exposes POST /score
   - Sometimes sleeps ~1.2s or returns 5xx (simulates issues)
   - Returns a risk score on success

2. **payments** (caller) - Runs on port 9000
   - Uses a Circuit Breaker (3 fails → Open 10s → 1 probe)
   - Calls Risk with 600ms timeout
   - If Open or call fails, returns fallback quickly

## Circuit Breaker States

- **Closed (normal)**: All calls pass through, counts consecutive failures
- **Open (tripped)**: Calls blocked immediately (fast-fail) for 10s cool-down
- **Half-Open (probe)**: After cool-down, allow one test call
  - Success → go Closed (resume normal)
  - Failure → go Open again (another cool-down)

## Parameters

- Trip threshold: 3 consecutive failures → Open
- Open duration (cool-down): 10 seconds
- Probe policy: 1 test call after 10s
- Caller timeout: 600ms per call

## How to Run

1. Start risk service:
```bash
cd risk
go run main.go
```

2. Start payments service:
```bash
cd payments
go run main.go
```

3. Test a payment:
```bash
curl -s -X POST http://localhost:9000/pay \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":1200}'
```

## Observations

- Normal responses when Risk is healthy
- Fallback responses when breaker is Open
- Small amounts (≤50) approved with challenge in fallback mode
- Large amounts held/declined in fallback mode