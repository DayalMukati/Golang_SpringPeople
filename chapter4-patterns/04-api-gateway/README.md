# API Gateway Pattern Example

This example demonstrates the API Gateway Pattern with routing, authentication, and rate limiting.

## Services

1. **gateway** - Runs on port 8000
   - Validates API keys
   - Applies rate limits (5 requests per 10 seconds per IP)
   - Adds request IDs for tracing
   - Routes requests to backend services

2. **users** - Runs on port 7001
   - Mock users service

3. **payments** - Runs on port 7002
   - Mock payments service

## How to Run

1. Start users service:
```bash
cd users
go run main.go
```

2. Start payments service:
```bash
cd payments
go run main.go
```

3. Start the gateway:
```bash
cd gateway
go run main.go
```

4. Call the gateway:
```bash
# Route to Users
curl -H "X-API-Key: demo-key-123" http://localhost:8000/users/123

# Route to Payments
curl -H "X-API-Key: demo-key-123" http://localhost:8000/payments/tx/abc
```