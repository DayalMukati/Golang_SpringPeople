# Introduction to Microservices Architecture

This example demonstrates basic microservices architecture with two independent services.

## Services

1. **user_service.go** - Runs on port 8081
   - Exposes GET /user
   - Returns user profile information

2. **payment_service.go** - Runs on port 8082
   - Exposes GET /pay
   - Returns payment confirmation

## How to Run

### Terminal 1 - User Service
```bash
go run user_service.go
```

### Terminal 2 - Payment Service
```bash
go run payment_service.go
```

### Test the Services

User Service:
```bash
curl http://localhost:8081/user
# Output: {"id":"U1001","name":"Ravi"}
```

Payment Service:
```bash
curl http://localhost:8082/pay
# Output: {"status":"Payment Successful","amount":100}
```

## Key Concepts

- Each service runs independently on different ports
- Services can be developed, deployed, and scaled separately
- Simulates a real microservices architecture