# 8. Capacity Planning and Resource Optimization

## Concepts

### What Is Capacity Planning?

**Capacity planning** is the process of predicting how much infrastructure (CPU, memory, storage, network) your system will need to handle current and future workloads without over-provisioning.

In cloud-native systems, this isn't just about "buying bigger servers." It's about:
- **Right-sizing pods** — deciding how much CPU/memory each container gets
- **Predicting workload spikes** — like salary days in fintech
- **Avoiding waste** — scaling down when demand is low

**Goal:** Meet user demand efficiently while minimizing cost.

### Why Is It Important?

**Without proper planning:**
- **Under-provisioning** → Transactions fail, apps slow down, customers lose trust
- **Over-provisioning** → Infrastructure bills skyrocket, burning business cash

**For fintech apps, both are unacceptable:**
- Users expect instant, reliable payments
- Businesses must control cloud costs

## FinPay Wallet Example

### Case Study: Salary Day

**Context:**
- Normal load = ~10,000 transactions/hour
- Salary day load = ~200,000 transactions/hour

**Without planning:**
- Run too few pods → system crashes
- OR always run 200,000-capacity → waste money 29 days a month

**Solution:**
1. Use historical data + monitoring to forecast spikes
2. Configure Horizontal Pod Autoscaler (HPA) to scale pods up at 70% CPU utilization
3. Use Cluster Autoscaler to add nodes if pods can't fit

**Result:** Smooth payments at midnight, without overspending during normal days

## Resource Optimization

Cloud-native platforms allow you to fine-tune resource requests and limits:

### Requests vs Limits

**Requests:**
- Minimum guaranteed CPU/memory for a pod
- Kubernetes scheduler uses this to place pods on nodes
- Pod gets **at least** this much

**Limits:**
- Maximum allowed CPU/memory a pod can consume
- If pod exceeds limit, it gets throttled (CPU) or killed (memory OOM)
- Pod gets **at most** this much

### Example in FinPay Wallet

**Payment API pod:**
- Request = 200m CPU, 256Mi memory
- Limit = 500m CPU, 512Mi memory

**Fraud Detection pod:**
- Request = 500m CPU, 512Mi memory
- Limit = 1 CPU, 1Gi memory

**Result:** Each service gets exactly what it needs. No pod starves; no pod hogs.

## Prerequisites

- Kubernetes cluster (minikube, kind, or Docker Desktop)
- kubectl
- metrics-server installed
- Basic understanding of YAML

## How to Run

### Step 1: Review Resource Configuration

```bash
cd chapter7-scalability/08-capacity-planning

# Review the deployment file
cat payment-deployment.yaml
```

**Key sections to understand:**

```yaml
resources:
  requests:
    cpu: 200m        # 0.2 CPU cores minimum
    memory: 256Mi    # 256 MiB minimum
  limits:
    cpu: 500m        # 0.5 CPU cores maximum
    memory: 512Mi    # 512 MiB maximum
```

### Step 2: Deploy with Resource Constraints

```bash
# Apply deployment
kubectl apply -f payment-deployment.yaml

# Check pod status
kubectl get pods

# View resource requests/limits
kubectl describe deployment payment-api
kubectl describe deployment fraud-detection
```

### Step 3: Monitor Resource Usage

```bash
# View current resource usage
kubectl top pods

# Watch resource usage in real-time
watch kubectl top pods

# View node capacity and usage
kubectl top nodes
```

**Expected output:**
```
NAME                              CPU(cores)   MEMORY(bytes)
payment-api-7d8f9c5b4-abc12      150m         200Mi
payment-api-7d8f9c5b4-def34      180m         220Mi
fraud-detection-6f4b8d9c-xyz89   450m         480Mi
```

### Step 4: Test Resource Limits

**Generate load to test CPU limits:**

```bash
# Port-forward to payment service
kubectl port-forward svc/payment-api 8080:80

# Generate load (in another terminal)
# Install hey: go install github.com/rakyll/hey@latest
hey -z 2m -c 100 http://localhost:8080/pay

# Watch CPU usage spike
watch kubectl top pods -l app=payment-api
```

**Observation:**
- CPU usage increases under load
- Stops at limit (500m) — Kubernetes throttles the pod
- Memory stays within limits

### Step 5: Analyze Resource Efficiency

```bash
# Get resource utilization percentage
kubectl describe node <node-name>

# Check for resource pressure
kubectl get pods --field-selector=status.phase=Pending
```

**Look for:**
- **High utilization (>80%)** → Consider increasing requests or adding nodes
- **Low utilization (<30%)** → Consider decreasing requests to save cost
- **Pending pods** → Not enough resources on nodes, need cluster autoscaler

### Step 6: Adjust Resources Based on Metrics

Based on monitoring data, you might update:

```yaml
# If payment API consistently uses 300m CPU
resources:
  requests:
    cpu: 250m        # Increase from 200m
    memory: 256Mi
  limits:
    cpu: 600m        # Increase from 500m
    memory: 512Mi
```

Then re-apply:
```bash
kubectl apply -f payment-deployment.yaml
kubectl rollout status deployment payment-api
```

## Benefits of Capacity Planning and Resource Optimization

### 1. Predictable Performance
- **Why it matters:** Customers expect financial transactions to be instant and reliable. A delayed payment erodes trust.
- **How planning helps:** By analyzing historical load patterns (salary days, festive shopping seasons), teams know when to expect spikes. Pods and databases are sized accordingly.

**Fintech Example:**
- FinPay Wallet sees highest activity on 1st of every month
- With proper capacity planning, system scales up before midnight
- Ravi pays his landlord at 12:05 AM and transfer is instant, despite 200,000 other users doing the same

**Outcome:** No slowdowns, no "try again later" errors — smooth experience even at peak times

### 2. Cost Efficiency
- **Why it matters:** Cloud bills grow quickly if every service is sized for peak load all the time
- **How planning helps:** Teams identify baseline usage vs peak usage, and configure autoscaling policies

**Fintech Example:**
- Instead of running 30 pods 24×7, FinPay Wallet runs 5 pods normally
- Scales up to 30 only during salary-day hours
- After rush, pods scale back down

**Outcome:** Business saves thousands in cloud costs each month while still delivering reliable service

### 3. Business Agility
- **Why it matters:** Fintech products evolve rapidly — new features, new customer segments, new regions create unpredictable demand
- **How planning helps:** With capacity planning, teams simulate expected loads before rolling out new features

**Fintech Example:**
- FinPay launches "Split Bills" during festival
- Forecasting shows likely 20% increase in transaction load
- Engineers adjust capacity settings before release

**Outcome:** Faster innovation with less firefighting

### 4. Operational Simplicity
- **Why it matters:** Without planning, engineers end up firefighting during every traffic spike
- **How planning helps:** Proper resource requests, limits, and autoscaling policies mean system self-adjusts

**Fintech Example:**
- On salary day, FinPay's autoscaler increases pods automatically
- SRE team doesn't need to wake up at midnight to manually add servers

**Outcome:** Stable operations, well-rested engineers, fewer human errors during critical windows

### 5. Better Resource Utilization
- **Why it matters:** Over-provisioned services waste CPU/memory, under-provisioned services cause slowdowns
- **How planning helps:** By measuring actual usage, teams assign right-sized requests and limits per service

**Fintech Example:**
- Fraud detection (CPU-heavy) given higher CPU limits
- Notification (I/O-heavy) optimized for memory
- Each service gets exactly what it needs

**Outcome:** No wasted resources, no starving services, better efficiency across cluster

## Capacity Planning Process

### Phase 1: Baseline Measurement (Week 1-2)

**Collect metrics:**
```bash
# CPU usage over time
kubectl top pods --all-namespaces

# Memory usage patterns
kubectl describe nodes

# Request latency (using application metrics)
curl http://prometheus:9090/api/v1/query?query=http_request_duration_seconds
```

**Questions to answer:**
- What's normal CPU/memory usage during off-peak?
- What's p95/p99 latency during normal load?
- How many requests per second during typical hours?

### Phase 2: Peak Load Analysis (Week 3-4)

**Analyze historical peaks:**
- Salary days (1st of month)
- Festival shopping days (Diwali, Christmas)
- Product launches

**Metrics to track:**
- Peak requests per second (RPS)
- Peak CPU/memory usage
- Error rate during peaks
- p95/p99 latency during peaks

### Phase 3: Capacity Modeling (Week 5)

**Calculate requirements:**

```
Peak RPS = 5,000 requests/second
Current capacity = 10 pods × 200 RPS = 2,000 RPS

Required pods = Peak RPS / RPS per pod
              = 5,000 / 200
              = 25 pods

Add 20% buffer = 25 × 1.2 = 30 pods
```

**CPU/Memory per pod:**
```
Observed peak CPU per pod = 400m
Add 20% headroom = 480m

Set:
  requests.cpu = 300m (scheduler placement)
  limits.cpu = 500m (throttle limit)
```

### Phase 4: Configure Autoscaling (Week 6)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: payment-hpa
spec:
  minReplicas: 5        # Baseline
  maxReplicas: 30       # Peak capacity
  targetCPUUtilization: 70%  # Scale when avg CPU > 70%
```

### Phase 5: Load Testing (Week 7)

```bash
# Simulate peak load
hey -z 10m -c 200 -q 50 http://payment-api/pay

# Monitor during test
watch kubectl get hpa
watch kubectl top pods
```

**Validate:**
- ✅ Pods scale up to handle load
- ✅ Latency stays below SLO (e.g., p95 < 500ms)
- ✅ Error rate stays below 0.1%
- ✅ Pods scale down after load drops

### Phase 6: Continuous Monitoring (Ongoing)

```bash
# Weekly capacity review
# - Check actual vs planned usage
# - Identify trending changes
# - Adjust requests/limits as needed

# Monthly load testing
# - Verify system still meets SLOs
# - Re-calibrate HPA thresholds
```

## Resource Optimization Techniques

### 1. Right-Sizing
- Start conservative, tune based on actual usage
- Monitor for 2-4 weeks before adjusting
- Aim for 60-80% utilization during normal load

### 2. Vertical Pod Autoscaler (VPA)
- Automatically adjusts requests/limits based on usage
- Use for recommendation, apply manually

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: payment-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: payment-api
  updateMode: "Off"  # Recommendation only
```

### 3. Quality of Service (QoS) Classes

**Guaranteed (highest priority):**
```yaml
requests:
  cpu: 500m
  memory: 512Mi
limits:
  cpu: 500m      # Same as requests
  memory: 512Mi  # Same as requests
```

**Burstable (medium priority):**
```yaml
requests:
  cpu: 200m
  memory: 256Mi
limits:
  cpu: 500m      # Higher than requests
  memory: 512Mi
```

**BestEffort (lowest priority):**
```yaml
# No requests or limits specified
# Can be evicted first under pressure
```

### 4. Node Affinity for Cost Optimization

```yaml
affinity:
  nodeAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      preference:
        matchExpressions:
        - key: node.kubernetes.io/instance-type
          operator: In
          values:
          - spot  # Use cheaper spot instances
```

## Capacity Planning Checklist

- [ ] Collect baseline metrics (CPU, memory, RPS) for 2-4 weeks
- [ ] Identify peak traffic patterns (time of day, day of week, month)
- [ ] Calculate required capacity with 20% buffer
- [ ] Set appropriate resource requests and limits
- [ ] Configure HPA with min/max replicas
- [ ] Enable Cluster Autoscaler for node scaling
- [ ] Perform load testing to validate
- [ ] Set up monitoring and alerting
- [ ] Document capacity model and assumptions
- [ ] Schedule monthly capacity reviews

## Troubleshooting

### Problem: Pods Pending

**Symptom:** Pods stuck in Pending state

**Diagnosis:**
```bash
kubectl describe pod <pod-name>
# Look for: "Insufficient cpu" or "Insufficient memory"
```

**Solutions:**
- Reduce resource requests
- Enable Cluster Autoscaler to add nodes
- Use smaller node instance types

### Problem: Pods OOMKilled

**Symptom:** Pods restarting with OOMKilled status

**Diagnosis:**
```bash
kubectl describe pod <pod-name>
# Look for: "Last State: Terminated, Reason: OOMKilled"
```

**Solutions:**
- Increase memory limits
- Fix memory leaks in application
- Add horizontal scaling to distribute load

### Problem: CPU Throttling

**Symptom:** High latency despite low CPU usage reported

**Diagnosis:**
```bash
# Check throttling metrics (requires cAdvisor)
kubectl exec -it <pod-name> -- cat /sys/fs/cgroup/cpu/cpu.stat
```

**Solutions:**
- Increase CPU limits
- Optimize application code
- Scale horizontally instead of vertically

## Key Takeaway

Capacity planning is about **balance:**
- Too little capacity → poor user experience
- Too much capacity → wasted money

For FinPay Wallet:
- **Normal days:** 5 pods with 200m CPU each = ₹500/month
- **Salary days:** 30 pods with 500m CPU each = ₹3,000/month (3 hours × 5 days)
- **Average monthly cost:** ~₹800/month

Compare to **always running peak capacity:**
- 30 pods × 500m CPU × 24×7 = ₹12,000/month

**Savings: 93% cost reduction** while maintaining performance

**Best practices:**
1. Monitor actual usage, not guesses
2. Start conservative, tune iteratively
3. Automate scaling with HPA/VPA
4. Load test before major events
5. Review capacity monthly

Capacity planning transforms cloud infrastructure from a **cost center** into a **competitive advantage** — delivering bank-grade performance at startup prices.