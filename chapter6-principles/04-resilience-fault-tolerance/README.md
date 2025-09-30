# Resilience and Fault Tolerance

This example demonstrates resilience through retry mechanisms and fault tolerance.

## Concepts

### Resilience
The ability of a system to recover quickly from disruptions and continue working.

### Fault Tolerance
The ability to keep operating correctly even when part of the system fails.

## Cloud-Native Approach

In cloud-native systems:
- Pods and nodes **will** fail (it's expected)
- Systems must be designed to handle failures gracefully
- Automatic recovery through retries, replicas, and self-healing
- Customers barely notice disruptions

## Files

1. **fraud.go** - Fraud service with simulated failures (30% failure rate)
2. **payment_with_retry.go** - Payment service with retry logic

## How It Works

### Scenario A: Fraud Service Intermittent Failure
1. Payment Service calls Fraud Detection
2. Fraud pod times out (simulated)
3. Payment Service retries after 500ms
4. Fraud responds with "SAFE" on second attempt
5. Payment completes successfully

From user's perspective: **No failure occurred**

### Scenario B: Kubernetes Auto-Healing
If a Fraud pod crashes:
1. Kubernetes detects unhealthy pod
2. Automatically spins up replacement pod
3. Load balancer routes to healthy pods
4. Service continues without downtime

## How to Run

### Terminal 1 - Fraud Service (with failures)
```bash
go run fraud.go
# Output: Fraud service (with simulated failures) running on :9090
```

### Terminal 2 - Payment Service (with retry)
```bash
go run payment_with_retry.go
# Output: Payment service with retry running on :8080
```

### Terminal 3 - Test Multiple Times
```bash
# Run multiple requests to see retry in action
for i in {1..5}; do
  curl http://localhost:8080/pay
  echo ""
  sleep 1
done
```

## Observe Resilience

Watch the terminals:
- **Fraud terminal**: Sometimes shows "Simulating fraud service failure..."
- **Payment terminal**: Shows "First attempt failed, retrying..."
- **Client**: Still gets successful response due to retry

## Code Explanation

### Simple Retry Logic
```go
resp, err := http.Get("http://localhost:9090/check")

if err != nil {
    // Retry once before failing
    time.Sleep(500 * time.Millisecond)
    resp, err = http.Get("http://localhost:9090/check")
    if err != nil {
        http.Error(w, "Fraud service unavailable", http.StatusServiceUnavailable)
        return
    }
}
```

**Why this works:**
- First attempt may fail if pod just restarted
- Most transient issues resolve quickly
- One retry significantly improves success rate
- Customer experience remains smooth

## Real-World Scenarios

### Scenario 1: Pod Deleted Accidentally
- Engineer accidentally deletes Wallet pod
- Kubernetes immediately spins up replacement
- Requests reroute to other healthy pods
- Payment still processes without downtime

### Scenario 2: Notification Service Down
- After wallet deducts money, Notification is unavailable
- Payment marks transaction as SUCCESS
- Notification queued to retry later
- Money transfer prioritized over alerts

### Scenario 3: Network Glitch
- Temporary network issue causes timeout
- Retry mechanism catches it
- Second attempt succeeds
- User never sees error

## Benefits

1. **Customer Trust**: Payments succeed even if one pod crashes
2. **Business Continuity**: Node failure doesn't stop transactions
3. **Graceful Degradation**: Non-critical services (Notification) can fail without blocking payments
4. **Operational Ease**: Kubernetes auto-restarts failed pods

## Kubernetes Fault Tolerance Features

### Automatic Pod Restart
```yaml
spec:
  replicas: 3  # Always maintain 3 replicas
  template:
    spec:
      restartPolicy: Always
```

### Health Checks
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 9090
  initialDelaySeconds: 5
  periodSeconds: 10
```

### Multiple Replicas
```bash
kubectl scale deployment fraud-deployment --replicas=5
```

## Advanced Patterns

1. **Circuit Breaker**: Stop calling failing service temporarily
2. **Bulkhead**: Isolate failures to prevent cascading
3. **Timeout**: Don't wait forever for unresponsive services
4. **Fallback**: Return cached/default values when service fails

## Fintech Context

In **FinPay Wallet**:
- **Resilience** ensures salary day transactions don't fail
- **Fault Tolerance** means one crashed pod doesn't stop millions of payments
- **Customer Trust** maintained through invisible recovery
- **Regulatory Compliance** met through reliable transaction processing

## Key Takeaway

In cloud-native systems, **failures are expected**.

The difference is:
- ❌ Traditional: One failure = entire system down
- ✅ Cloud-Native: Failures handled gracefully, system keeps running