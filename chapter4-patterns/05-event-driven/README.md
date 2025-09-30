# Event-Driven Communication Example

This example demonstrates Event-Driven Communication between a payments producer and a ledger consumer.

## Services

1. **payments** (producer) - Runs on port 9000
   - Exposes POST /authorize
   - Authorizes payment requests
   - Emits PaymentAuthorized events asynchronously

2. **ledger** (consumer) - Runs on port 9001
   - Exposes POST /events
   - Receives PaymentAuthorized events
   - Converts to ledger entries (double-entry demo)
   - Prints entries (simulates persistence)

## How to Run

1. Start the ledger (consumer):
```bash
cd ledger
go run main.go
```

2. Start the payments (producer):
```bash
cd payments
go run main.go
```

3. Send a test request:
```bash
curl -X POST http://localhost:9000/authorize \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":42.50,"currency":"USD","merchant_id":"M10"}'
```

4. Check the ledger terminal for the printed ledger entry.