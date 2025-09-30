# Sidecar Pattern Example

This example demonstrates the Sidecar Pattern with a payments service and a logging sidecar.

## Services

1. **payments** (main app) - Runs on port 8080
   - Exposes POST /authorize
   - Processes payment authorization requests
   - Sends audit events to the sidecar

2. **logsidecar** (helper) - Runs on port 9000
   - Exposes POST /logs
   - Receives audit events from payments
   - Enriches events with timestamp and fraud score
   - Prints enriched records (simulates storage)

## How to Run

1. Start the sidecar:
```bash
cd logsidecar
go run main.go
```

2. Start the payments app (in another terminal):
```bash
cd payments
go run main.go
```

3. Send a test request:
```bash
curl -X POST http://localhost:8080/authorize \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":1200,"currency":"USD","merchant_id":"M10","ip_address":"203.0.113.10","device_id":"D-abc"}'
```

4. Check the sidecar terminal for the enriched log record.