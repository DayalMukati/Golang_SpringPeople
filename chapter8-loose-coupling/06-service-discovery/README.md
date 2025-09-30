# 6. Service Discovery and Dynamic Binding

## Concepts

### What Is Service Discovery?

In traditional applications, services usually live on fixed servers with fixed IP addresses. A payment system might call Fraud at `http://10.0.0.5:8080` and Notification at `http://10.0.0.6:9090`.

**But in cloud-native systems, things change:**
- Containers start and stop dynamically
- Pods are replaced automatically when they fail
- Scaling creates many instances of the same service
- Each pod gets a new IP every time it starts

**Hardcoding addresses no longer works.** Instead, systems need **service discovery** — a mechanism that lets services find each other dynamically without knowing exact IPs.

### What Problem Does It Solve?

**1. Dynamic IPs**
- Every time Kubernetes restarts a pod, it gets a new IP
- If Payment had a hardcoded Fraud IP, it would break instantly

**2. Scaling**
- Fraud may run on 10 pods today, 15 tomorrow
- Payment must load-balance across them automatically

**3. Resilience**
- If one Fraud pod dies, Payment should automatically connect to another healthy one

**Service discovery removes the need for manual configuration and ensures services always find each other in a changing environment.**

## Real-World Example: FinPay Wallet

### Scenario: Ravi Makes a Rent Transfer

**Requirements:**
- Payment Service needs to talk to Fraud Service
- But Fraud is running on 10 pods, each with a different IP
- Tomorrow, Kubernetes may restart them with new IPs

**With service discovery:**
- Payment doesn't care about pod addresses
- It just calls: `http://fraud-service:8080/check`
- Kubernetes' DNS-based service discovery automatically resolves this name to the available Fraud pods

**Result:** Even if pods restart or scale up/down, Payment still works without code changes

## How It Works in Kubernetes

### 1. DNS Service Names
- Each service (like Fraud) is given a stable name (`fraud-service`)
- Payment uses this name, not the pod IP

### 2. Load Balancing
- Behind the service name, Kubernetes load-balances requests across all Fraud pods

### 3. Service Mesh (Advanced)
- Tools like Istio or Linkerd add dynamic routing, retries, and failover
- **Example:** If one Fraud pod is slow, the mesh routes traffic to a healthier pod

## Prerequisites

- Go 1.20+
- Understanding of environment variables
- Basic understanding of Kubernetes services (optional)

## How to Run

### Step 1: Run with Default Service URL

```bash
cd chapter8-loose-coupling/06-service-discovery

# Run without environment variable (uses default)
go run discovery_demo.go
```

**Expected output:**
```
Payment Service discovered Fraud at: http://fraud-service:8080
```

### Step 2: Run with Custom Service URL

```bash
# Set environment variable
export FRAUD_SERVICE_URL="http://fraud-service.finpay.svc.cluster.local:8080"

# Run again
go run discovery_demo.go
```

**Expected output:**
```
Payment Service discovered Fraud at: http://fraud-service.finpay.svc.cluster.local:8080
```

### Step 3: Understand the Code

```go
func main() {
    // Simulate discovery: Fraud service address comes from config/env
    fraudService := os.Getenv("FRAUD_SERVICE_URL")

    if fraudService == "" {
        fraudService = "http://fraud-service:8080" // default
    }

    fmt.Println("Payment Service discovered Fraud at:", fraudService)
    // Normally here we would call fraudService via HTTP
}
```

**Key points:**
- Payment doesn't hardcode Fraud's IP address
- It reads from environment variable (like Kubernetes does)
- In real Kubernetes, `FRAUD_SERVICE_URL` would be automatically injected as `fraud-service:8080`
- Payment doesn't care about Fraud pod IPs — it just uses the stable service name

**This is the essence of dynamic binding.**

## Kubernetes Service Discovery

### Without Service Discovery (Hardcoded IPs)

```go
// ❌ Bad: Hardcoded IP addresses
func callFraudService() {
    resp, _ := http.Get("http://10.0.0.5:8080/check")
    // ...
}
```

**Problems:**
- Pod restarts → new IP → code breaks
- Scaling to 10 pods → must manually update all 10 IPs
- Pod failure → must manually route to another IP

### With Service Discovery (DNS Names)

```go
// ✅ Good: Use service name
func callFraudService() {
    resp, _ := http.Get("http://fraud-service:8080/check")
    // ...
}
```

**Benefits:**
- Pod restarts → same service name works
- Scaling to 10 pods → Kubernetes load-balances automatically
- Pod failure → Kubernetes routes to healthy pods automatically

## Kubernetes DNS Resolution

### Service Name Format

**Within same namespace:**
```
http://fraud-service:8080
```

**Cross-namespace:**
```
http://fraud-service.finpay-prod.svc.cluster.local:8080
```

**Format breakdown:**
```
<service-name>.<namespace>.svc.cluster.local
```

### Example Kubernetes Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: fraud-service
  namespace: finpay-prod
spec:
  selector:
    app: fraud
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

**How it works:**
1. Kubernetes creates DNS entry: `fraud-service.finpay-prod.svc.cluster.local`
2. Payment calls this DNS name
3. Kubernetes DNS resolves to one of the healthy Fraud pod IPs
4. Load balances across all pods matching selector `app: fraud`

## Real-World Flow (FinPay Example)

### Scenario: Salary Day with Auto-Scaling

**Initial state:**
- Fraud Service: 5 pods
- Pod IPs: 10.0.0.10, 10.0.0.11, 10.0.0.12, 10.0.0.13, 10.0.0.14

**Payment calls:**
```go
http.Get("http://fraud-service:8080/check")
```

**Kubernetes resolves:**
- DNS lookup: `fraud-service` → 5 pod IPs
- Load balances request to one pod
- Payment succeeds ✅

**Salary day surge → HPA scales Fraud to 15 pods:**
- New pod IPs: 10.0.0.15 through 10.0.0.24
- **Payment code unchanged** — still calls `fraud-service:8080`
- Kubernetes automatically includes new pods in load balancing

**Pod 10.0.0.12 crashes:**
- Kubernetes removes it from service endpoints
- **Payment code unchanged** — DNS resolution excludes crashed pod
- Requests route to remaining 14 healthy pods

**Result:**
- Zero code changes in Payment Service
- Zero downtime
- Automatic scaling and failover

## Service Discovery Patterns

### 1. Client-Side Discovery (Less Common)

```
Payment → Service Registry (etcd/Consul) → Gets list of IPs → Calls Fraud
```

**Pros:** Client controls load balancing
**Cons:** More complexity in client code

### 2. Server-Side Discovery (Kubernetes Default)

```
Payment → Kubernetes Service → Load balances → Fraud pods
```

**Pros:** Simple client code, Kubernetes handles everything
**Cons:** Less control over load balancing

### 3. Service Mesh (Advanced)

```
Payment → Envoy Sidecar → Service Mesh → Fraud pods
```

**Pros:** Advanced routing, retries, circuit breakers, observability
**Cons:** Additional complexity and resource overhead

## Benefits of Service Discovery

### 1. Zero Configuration
- No manual IP management
- No configuration file updates when pods restart

### 2. Automatic Load Balancing
- Kubernetes distributes requests across all healthy pods
- No manual load balancer configuration

### 3. Dynamic Scaling
- Add pods → automatically included in load balancing
- Remove pods → automatically excluded

### 4. Fault Tolerance
- Unhealthy pods automatically removed from rotation
- Requests route to healthy pods only

### 5. Multi-Environment Support
- Same code works in dev, staging, prod
- Different service endpoints per environment (via DNS)

## Environment Variable Configuration

### Development

```bash
export FRAUD_SERVICE_URL="http://localhost:9090"
```

### Kubernetes (Staging)

```bash
export FRAUD_SERVICE_URL="http://fraud-service.finpay-staging.svc.cluster.local:8080"
```

### Kubernetes (Production)

```bash
export FRAUD_SERVICE_URL="http://fraud-service.finpay-prod.svc.cluster.local:8080"
```

**Same code, different environments — powered by service discovery**

## Key Takeaway

**Service discovery is essential for cloud-native systems where services are dynamic and ephemeral.**

For FinPay Wallet:
- **Without service discovery:** Payment breaks every time Fraud pods restart (new IPs)
- **With service discovery:** Payment always finds Fraud using stable DNS name, regardless of pod IPs

**Design principle:** "Don't depend on IPs, depend on names."

In Kubernetes:
- Services provide stable DNS names
- Pods can come and go freely
- Load balancing happens automatically
- Failover is built-in

**"Service discovery transforms rigid, IP-based systems into flexible, cloud-native architectures where services find each other dynamically."**