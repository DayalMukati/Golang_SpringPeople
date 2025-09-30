# Loose Coupling Between Services

This example demonstrates loose coupling through independent services communicating via APIs.

## Concept

**Loose Coupling** means services work independently:
- Each service has its own responsibility
- Services communicate through well-defined APIs
- A change in one service doesn't break others
- Services don't share databases or hidden dependencies

## Services

1. **fraud.go** (Fraud Detection Service)
   - Runs on port 9090
   - Exposes GET /check endpoint
   - Returns fraud check result ("SAFE")
   - Independent service with its own logic

2. **payment.go** (Payment Service)
   - Runs on port 8080
   - Exposes GET /pay endpoint
   - Calls Fraud Service via HTTP
   - Doesn't know Fraud's internal implementation

## Why This Is Loose Coupling

- **Independent deployment**: Fraud team can update ML models without touching Payment code
- **API contract**: Payment only knows the endpoint (`/check`), not the logic
- **Separate databases**: Each service manages its own data
- **Failure isolation**: If Fraud crashes, Payment can handle it gracefully
- **Team autonomy**: Different teams own different services

## How to Run

### Terminal 1 - Fraud Service
```bash
go run fraud.go
# Output: Fraud service running on :9090
```

### Terminal 2 - Payment Service
```bash
go run payment.go
# Output: Payment service running on :8080
```

### Terminal 3 - Test
```bash
curl http://localhost:8080/pay
# Output: Payment processed. Fraud check result: SAFE
```

## Real-World Flow

When a user pays ₹5,000:

1. **Payment Service** receives request
2. Calls Fraud Detection: `GET http://fraud-service:9090/check`
3. **Fraud Service** returns: `{ "status": "SAFE" }`
4. **Payment Service** proceeds with transaction
5. Calls **Wallet Service**: `POST http://wallet-service:6060/debit`
6. Calls **Notification Service**: `POST http://notification-service:7070/send`

Each service is independent, communicates via APIs, and can be deployed separately.

## Benefits

1. **Independent updates**: Upgrade Fraud without touching Payments
2. **Parallel development**: Teams work simultaneously
3. **Failure isolation**: One service down doesn't crash others
4. **Technology freedom**: Each service can use different languages/databases
5. **Easy testing**: Services can be tested in isolation

## Contrast: Tightly Coupled

In a tightly coupled system:
- All logic in one monolithic application
- Shared database causes data coupling
- Bug in Notification crashes entire app
- Cannot scale services independently
- Deployment requires entire app rebuild

## Fintech Example

**FinPay Wallet** uses loose coupling:
- Payment Service handles transactions
- Fraud Service runs ML models (Python)
- Wallet Service manages balances (Go)
- Notification Service sends alerts (Node.js)

Each service evolves independently while working together seamlessly.