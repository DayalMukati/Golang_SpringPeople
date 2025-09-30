# Benefits and Challenges of Microservices

This example demonstrates one of the challenges: network latency between microservices.

## Challenge Demonstrated

The payment service simulates calling a Fraud Service with a 2-second delay.

## How to Run

```bash
go run payment_service.go
```

## Test

```bash
curl http://localhost:8080/pay
# Takes 2 seconds to respond
# Output: Payment processed successfully
```

## Key Concepts

### Benefits
- Independent scaling
- Faster development and deployment
- Resilience and fault isolation
- Technology diversity
- Clear team ownership

### Challenges
- **Network overhead** (demonstrated here with time.Sleep)
- Operational complexity
- Data management
- Debugging and monitoring
- Testing complexity

## Solution

In production, this would be solved with:
- Async processing
- Circuit breakers
- Better infrastructure scaling
- Caching