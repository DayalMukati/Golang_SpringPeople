# 1. Introduction to Scalability

## Concepts

**Scalability** means a system can grow with demand. If tomorrow your app gets ten times more users, will it continue to perform well, or will it collapse under the weight of new requests?

In a cloud-native world, scalability isn't just about adding more power; it's about **elasticity** — the ability to grow when traffic spikes and shrink when things quiet down.

### Why Scalability Matters

Applications that ignore scalability eventually hit limits:
- **Slower response times**: The app feels laggy when requests pile up
- **Downtime during peaks**: A sudden burst of activity crashes the system
- **Inefficient costs**: Over-provisioned servers sit unused most of the time

In fintech especially, scalability is critical. When people are dealing with money, they expect systems to "just work."

### FinPay Wallet Example

**Without scalability:**
- Regular day: 10,000 transactions/hour ✓
- Salary day: 200,000 transactions/hour → System overload
- Payment service instances get overloaded
- Transactions queue, timeout, or fail
- Customers frustrated when transfers don't go through

**With scalability:**
- Kubernetes detects the surge in demand
- Automatically spins up more pods (3 → 30 payment service pods)
- Load balancer distributes traffic across them
- Payments processed instantly, just like any normal day

### Key Lesson

**Scalability is achieved more through deployment architecture than through application code.**

## Prerequisites

- Go 1.20+
- Basic understanding of HTTP servers
- Docker (optional, for containerization)
- Kubernetes (optional, for orchestration)

## How to Run

### Local Development

```bash
cd chapter7-scalability/01-introduction

# Run the service
go run main.go

# Test the endpoint
curl http://localhost:8080/pay
```

**Expected output:**
```
Payment processed
```

### Understanding the Code

This simple Go service handles payments, but only as much as the server it runs on allows. On one laptop, maybe a few hundred requests per second.

**To scale it:**
1. We don't rewrite the code
2. We run multiple instances behind a load balancer
3. Kubernetes/Docker Swarm can launch 10, 20, or 30 instances depending on traffic

### Build Docker Image (Optional)

```bash
# Create Dockerfile
cat > Dockerfile <<EOF
FROM golang:1.21 AS builder
WORKDIR /app
COPY main.go .
RUN go mod init payment && CGO_ENABLED=0 go build -o payment main.go

FROM debian:stable-slim
COPY --from=builder /app/payment /payment
EXPOSE 8080
ENTRYPOINT ["/payment"]
EOF

# Build and run
docker build -t payment:v1 .
docker run -p 8080:8080 payment:v1
```

### Deploy to Kubernetes (Optional)

```bash
# Create deployment with 3 replicas
kubectl create deployment payment --image=payment:v1 --replicas=3

# Expose as a service
kubectl expose deployment payment --port=80 --target-port=8080

# Scale up during peak
kubectl scale deployment payment --replicas=30

# Scale down after peak
kubectl scale deployment payment --replicas=3
```

## Benefits of Scalability

### 1. Consistent User Experience
- Performance remains stable even during peak traffic
- Example: On salary day, millions log into FinPay simultaneously
- With scaling, Ravi's ₹5,000 rent payment goes through as quickly as on a quiet Tuesday

### 2. Support for Business Growth
- System grows organically with demand
- No need to rebuild the app when new users join
- Example: FinPay grows from 100K to 1M users after bank partnership
- Simply increase pods and database partitions

### 3. Cost Efficiency
- Handle users efficiently, not just handle more users
- Example: Run 30 pods on salary day, scale down to 5 the next day
- Pay for extra compute only when needed

### 4. Reliability
- Systems aren't tied to a single instance
- Example: If 1 pod crashes out of 30, load balancer reroutes to remaining 29
- Users never notice the failure

## Real-World Flow (FinPay Example)

### Normal Day
- Traffic: 10,000 transactions/hour
- Pods: 3 payment service instances
- CPU: ~40% utilization
- Status: ✓ All systems normal

### Salary Day Midnight
- Traffic: 200,000 transactions/hour (20× spike)
- Kubernetes HPA detects high CPU (>80%)
- Pods: Auto-scales 3 → 10 → 20 → 30
- Load balancer distributes traffic evenly
- CPU: ~70% utilization across all pods
- Status: ✓ All payments processing smoothly

### Next Day Afternoon
- Traffic: Returns to 12,000 transactions/hour
- Kubernetes scales down: 30 → 10 → 5 pods
- Cost: Reduced cloud spend
- Status: ✓ Efficient resource usage

## Key Takeaway

Scalability is not just a technical feature — it directly impacts customer trust and business reputation. In fintech, where every transaction matters, the ability to scale seamlessly during peak demand while controlling costs during quiet periods is essential for success.