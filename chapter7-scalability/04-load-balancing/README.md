# 4. Load Balancing Strategies

## Concepts

**Load balancing** distributes incoming requests across multiple instances of your service (pods/VMs) so no single instance becomes a bottleneck.

### Two Layers of Load Balancing

**Layer 4 (TCP/UDP):**
- Balances connections based on IP/port
- Fast, simple, protocol-agnostic
- Use case: Generic traffic distribution

**Layer 7 (HTTP/HTTPS/gRPC):**
- Understands requests (paths, headers, cookies)
- Can route `/api/pay` to one service and `/api/refund` to another
- Can terminate TLS, enforce auth, etc.
- Use case: Smart routing, canary deployments, A/B testing

### What Problem Does It Solve?

**Without a load balancer:**
- Traffic piles onto whichever instance clients happen to hit
- One busy instance slows down while others sit idle
- Failed instances continue receiving requests

**With a proper load balancer:**
- Shares load so latency stays low
- Masks failures by skipping unhealthy instances
- Enables scale-out (add pods, balancer starts sending them work)
- Supports smarter routing (version canaries, geo routing, path rules)

## Fintech Example: FinPay Wallet

FinPay's Payment service runs 3 pods on normal days and 30 on salary day. Customers fire `/pay` requests from mobile apps.

**How load balancing helps:**
1. A Kubernetes Service (ClusterIP) fronts the pods
2. kube-proxy sends connections to healthy pods
3. An external LoadBalancer/Ingress receives public HTTPS, terminates TLS, forwards to Service
4. During rush: new pods come online, balancer discovers them and distributes traffic automatically
5. If Pod 7 crashes mid-rush, health checks fail and balancer stops sending it traffic
6. **Customers never notice the failure**

## Common Load Balancing Algorithms

### 1. Round Robin
**How it works:** Sends requests to instances in order (Pod1 → Pod2 → Pod3 → Pod1 → ...)

**Use when:** Pods are similar size and requests are short and uniform (typical JSON APIs)

**Example:** FinPay payment authorizations — each request takes ~50ms

### 2. Least Connections
**How it works:** Picks the instance with the fewest active connections

**Use when:** Request lengths vary a lot

**Example:** Some fraud checks are heavier than others (simple checks: 10ms, complex ML checks: 500ms)

### 3. Weighted Round Robin
**How it works:** Give bigger/faster pods more weight

**Use when:** You mix instance sizes

**Example:** Some nodes are larger (8-core vs 4-core), give them 2× traffic

### 4. Consistent Hashing (Sticky Routing)
**How it works:** Hash on a key (userId, walletId) so same user tends to hit same pod

**Use when:** You keep short-lived, in-memory state like per-user rate limits or small cache

**Caution:** Don't depend on it for critical state — persist to DB/cache cluster

### 5. Random (with Two Choices)
**How it works:** Pick two pods at random, choose the less loaded

**Benefits:** Surprisingly effective at scale, simple to implement

## Load Balancing in Kubernetes

### Service (ClusterIP/NodePort/LoadBalancer)
- Provides stable virtual IP + endpoint set
- kube-proxy (iptables/ipvs) spreads connections across ready pods
- **Readiness probes matter!** Only ready pods receive traffic

### Ingress (Layer 7)
Ingress Controllers (Nginx, HAProxy, Envoy, Traefik) provide:
- HTTP-aware routing: host/path rules
- TLS termination
- Sticky sessions
- Canary headers

### Service Mesh (Optional)
Sidecars (Envoy) add per-request:
- Load balancing
- Timeouts
- Retries
- Circuit breakers
- Traffic splitting

**All without changing app code**

## Prerequisites

- Go 1.20+
- Kubernetes cluster (minikube, kind, or Docker Desktop)
- kubectl
- Docker Hub account

## How to Run

### Step 1: Build and Push Image

```bash
cd chapter7-scalability/04-load-balancing

# Initialize Go module
go mod init loadbalancer-demo

# Test locally
go run main.go

# In another terminal
curl http://localhost:8080/pay

# Build Docker image
cat > Dockerfile <<EOF
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o payment main.go

FROM debian:stable-slim
COPY --from=builder /app/payment /payment
EXPOSE 8080
ENTRYPOINT ["/payment"]
EOF

# Build and push
docker build -t <YOUR_DOCKERHUB_USERNAME>/payment-lb:v1 .
docker push <YOUR_DOCKERHUB_USERNAME>/payment-lb:v1
```

### Step 2: Deploy to Kubernetes

```bash
# Create deployment with 3 replicas
kubectl create deployment payment-lb --image=<YOUR_DOCKERHUB_USERNAME>/payment-lb:v1 --replicas=3

# Expose as ClusterIP service
kubectl expose deployment payment-lb --port=80 --target-port=8080 --name=payment-lb-service

# Check pods
kubectl get pods -l app=payment-lb
```

### Step 3: Test Load Balancing (Round Robin)

**Option A: Port-forward and test**
```bash
# Port-forward to service
kubectl port-forward svc/payment-lb-service 8080:80

# In another terminal, curl multiple times
for i in {1..10}; do
  curl -s http://localhost:8080/pay | jq .pod
done
```

**Expected output:** Pod names alternate (Round Robin)
```
payment-lb-7d8f9c5b4-abc12
payment-lb-7d8f9c5b4-def34
payment-lb-7d8f9c5b4-ghi56
payment-lb-7d8f9c5b4-abc12
payment-lb-7d8f9c5b4-def34
...
```

**Option B: Debug pod inside cluster**
```bash
# Create debug pod
kubectl run curl-debug --image=curlimages/curl -i --tty --rm -- sh

# Inside the pod
for i in 1 2 3 4 5 6 7 8 9 10; do
  curl -s http://payment-lb-service/pay | grep pod
done
```

### Step 4: Simulate Pod Failure

```bash
# Delete one pod
kubectl delete pod <pod-name>

# Immediately test
for i in {1..20}; do
  curl -s http://localhost:8080/pay | jq .pod
  sleep 0.2
done
```

**Observation:**
- Failed pod stops receiving traffic immediately (readiness probe fails)
- Other pods continue serving requests
- Kubernetes replaces failed pod automatically
- No user-facing errors

### Step 5: Scale and Observe Distribution

```bash
# Scale to 6 replicas
kubectl scale deployment payment-lb --replicas=6

# Wait for pods to be ready
kubectl get pods -l app=payment-lb -w

# Test again
for i in {1..20}; do
  curl -s http://localhost:8080/pay | jq .pod
done
```

**Observation:** Traffic now distributes across 6 pods

## Health Checks, Timeouts, and Retries

### Why They Matter

**Readiness probe:** Only ready pods receive traffic
```yaml
readinessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

**Liveness probe:** Restarts stuck pods
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
```

**Timeouts & retries:** Balancer or client library should give up quickly on slow pods and retry another (idempotent operations only)

### Important for Fintech

In a payment flow:
- Keep authorization checks idempotent so safe retries don't double-charge
- For non-idempotent operations, use request IDs and server-side deduplication
- Set sane timeouts (300-800ms per hop)
- Never use infinite retries

## Benefits of Load Balancing

### 1. Even Distribution
- No single pod becomes a bottleneck
- All pods share the load equally
- **Example:** 30 pods each handle ~6,600 requests/minute instead of one pod handling 200,000

### 2. High Availability
- Failed pods are automatically removed from rotation
- **Example:** Pod 7 crashes, balancer routes traffic to remaining 29 pods
- Users experience 0 downtime

### 3. Seamless Scaling
- New pods automatically receive traffic
- **Example:** HPA scales from 3 to 10 pods, load balancer discovers new pods via Kubernetes endpoints API
- No manual configuration needed

### 4. Health-Based Routing
- Only healthy pods receive traffic
- **Example:** Pod 5 becomes slow (high memory usage), readiness probe fails, traffic stops
- System self-heals by isolating problematic pods

### 5. Advanced Routing (Layer 7)
- Path-based routing: `/api/pay` → payment service, `/api/fraud` → fraud service
- Canary deployments: 5% traffic to v2, 95% to v1
- A/B testing: Users from region X → new feature, others → old feature

## Practical Trade-offs and Tips

### Start Simple
- Round Robin + solid readiness probes solves most problems
- Don't over-engineer until you have specific needs

### Sticky Sessions (Use Sparingly)
- Hurts even distribution
- Prefer externalizing state to Redis/DB
- Only use for short-lived session data

### Guard Rails
- Set sane timeouts (300-800ms per hop)
- Use retries with jitter
- Never infinite retries

### Observability
- Export per-pod request counters
- Monitor distribution: imbalanced pattern usually means readiness/annotation issue or slow pod

### Canaries & Blue/Green
- L7 balancers can route small % of traffic to new version
- Roll forward if healthy, roll back fast if not
- **Example:** Deploy payment-v2, route 5% traffic, monitor error rate, gradually increase to 100%

## Real-World Flow (FinPay Example)

### Scenario: Salary Day with Load Balancing

**00:00 - Surge begins**
- Traffic: 200,000 requests/minute
- Initial pods: 3
- Load balancer: Round-robin across 3 pods
- Each pod: ~66,666 requests/minute → CPU spikes to 95%

**00:02 - Auto-scaling kicks in**
- HPA scales to 10 pods
- Load balancer automatically discovers new pods
- Each pod: ~20,000 requests/minute → CPU drops to 75%

**00:05 - Peak load**
- HPA scales to 20 pods
- Each pod: ~10,000 requests/minute → CPU stabilizes at 70%
- All requests processed smoothly

**00:10 - Pod 7 crashes (bug)**
- Readiness probe fails
- Load balancer stops sending traffic to Pod 7
- Remaining 19 pods absorb traffic
- Each pod: ~10,526 requests/minute → CPU rises slightly to 72%
- **Users notice nothing**
- Kubernetes replaces Pod 7 automatically

**14:00 - Traffic normalizes**
- HPA scales down to 5 pods
- Load balancer adjusts
- Each pod: ~2,400 requests/minute → CPU at 40%
- Cost-optimized

## Key Takeaway

Load balancing is the foundation of horizontal scalability. Without it, you can't effectively distribute traffic across multiple instances.

For fintech applications like FinPay Wallet, load balancing ensures:
- **High availability** — failed pods don't impact users
- **Even performance** — no pod becomes a bottleneck
- **Seamless scaling** — add capacity without downtime
- **Operational simplicity** — automated health-based routing

Combined with auto-scaling, load balancing enables cloud-native fintech systems to handle unpredictable traffic patterns reliably and cost-effectively.