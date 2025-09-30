# 3. Auto-Scaling in Cloud Environments

## Concepts

**Auto-scaling** is the ability of a cloud system to automatically increase or decrease resources depending on workload. Instead of engineers watching dashboards and adding servers manually, the platform continuously monitors metrics (CPU, memory, or custom signals like request rate) and adjusts capacity on its own.

### Kubernetes Auto-Scaling Components

**1. Horizontal Pod Autoscaler (HPA)**
- Adjusts the number of pods for a Deployment based on metrics like CPU or memory utilization
- Example: Scale from 3 to 30 pods when CPU > 80%

**2. Cluster Autoscaler (CA)**
- If pods can't be scheduled because there's no room on worker nodes, CA adds more nodes
- When demand falls, it removes extra nodes

Together, these make the system **elastic**: growing during spikes, shrinking when idle.

### What Problem Does It Solve?

**Without auto-scaling:**
- Teams must guess peak load
- Either over-provision (waste money during off-hours) or under-provision (system crashes during peaks)

**With auto-scaling:**
- Performance stays steady under sudden traffic bursts
- Costs remain efficient when demand falls
- No human intervention needed at midnight during peak traffic

## FinPay Wallet Example

### Scenario: Salary Day Traffic Surge

**Normal weekday:**
- ~10,000 transactions per hour
- 3 Payment Service pods running
- CPU usage: ~40%

**Salary day midnight:**
- ~200,000 transactions per hour flood in
- Payment Service pods see CPU usage >80%
- Kubernetes HPA triggers and scales pods gradually: 3 → 10 → 20 → 30
- Load Balancer routes traffic evenly across new pods
- Customers complete payments instantly — they never notice scaling behind the scenes

**After Peak (afternoon):**
- Traffic falls back to normal
- HPA scales pods back down to 5
- Cloud costs reduced

**Result:** Auto-scaling protects user experience during spikes and saves money when load drops.

## Prerequisites

- Go 1.20+
- Kubernetes cluster with metrics-server installed
- kubectl
- Docker Hub account (for custom images)

## How to Run

### Step 0: Verify Metrics Server

```bash
# Check if metrics-server is running
kubectl get deployment metrics-server -n kube-system

# If not installed (for minikube)
minikube addons enable metrics-server

# For other clusters, install metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### Step 1: Create Payment Service Deployment

```bash
cd chapter7-scalability/03-auto-scaling

# Create deployment (use image from chapter7-scalability/01-introduction)
kubectl create deployment payment-deployment --image=<YOUR_DOCKERHUB_USERNAME>/payment:v1 --replicas=3

# Set resource requests (required for HPA)
kubectl set resources deployment payment-deployment --requests=cpu=100m,memory=128Mi --limits=cpu=500m,memory=512Mi

# Expose as service
kubectl expose deployment payment-deployment --port=80 --target-port=8080 --name=payment-service
```

### Step 2: Apply HPA

```bash
# Apply the HPA configuration
kubectl apply -f payment-hpa.yaml

# Check HPA status
kubectl get hpa payment-hpa

# Watch HPA in action
kubectl get hpa payment-hpa -w
```

### Step 3: Generate Load to Trigger Scaling

```bash
# Get service URL
# For minikube
export SERVICE_URL=$(minikube service payment-service --url)

# For kind or Docker Desktop
kubectl port-forward svc/payment-service 8080:80
export SERVICE_URL=http://localhost:8080

# Generate load using a simple loop
while true; do curl -s $SERVICE_URL/pay > /dev/null; done

# Or use a load testing tool like 'hey'
# Install: go install github.com/rakyll/hey@latest
hey -z 5m -c 50 $SERVICE_URL/pay
```

### Step 4: Observe Auto-Scaling

**Terminal 1: Watch HPA**
```bash
kubectl get hpa payment-hpa -w
```

**Expected output:**
```
NAME          REFERENCE                       TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
payment-hpa   Deployment/payment-deployment   45%/80%   3         30        3          1m
payment-hpa   Deployment/payment-deployment   85%/80%   3         30        3          2m
payment-hpa   Deployment/payment-deployment   92%/80%   3         30        6          3m
payment-hpa   Deployment/payment-deployment   78%/80%   3         30        10         4m
```

**Terminal 2: Watch Pods**
```bash
kubectl get pods -l app=payment-deployment -w
```

**Observation:** Pods increase from 3 → 6 → 10 → ... as CPU stays above 80%

### Step 5: Stop Load and Watch Scale-Down

```bash
# Stop the load generator (Ctrl+C)

# Watch HPA scale down (takes 5-10 minutes)
kubectl get hpa payment-hpa -w
```

**Expected behavior:**
- CPU drops below 80%
- After cooldown period (~5 minutes), HPA scales pods back down
- Eventually returns to minReplicas: 3

### Step 6: Check Resource Utilization

```bash
# View current CPU/memory usage
kubectl top pods -l app=payment-deployment

# View node utilization
kubectl top nodes
```

## Understanding the HPA Configuration

```yaml
minReplicas: 3          # Always keep at least 3 pods
maxReplicas: 30         # Don't go beyond 30 (protects from runaway scaling)
averageUtilization: 80% # If pods use more than 80% CPU, add more pods
```

**How it works:**
1. Metrics server collects CPU/memory usage from all pods every 15 seconds
2. HPA calculates average utilization across all pods
3. If average > 80%, HPA increases replicas
4. If average < 80% for sustained period, HPA decreases replicas
5. Scaling happens gradually to avoid thrashing

## Benefits of Auto-Scaling

### 1. Hands-Free Operations
- Engineers don't need to manually add servers at midnight
- **Example:** FinPay engineers sleep while HPA scales up pods automatically during salary day surge

### 2. Performance Stability
- System adapts in real time to handle traffic spikes
- **Example:** No transaction delays even during 20× traffic spikes
- Customers experience consistent response times

### 3. Cost Efficiency
- Pay for what you use, not for worst-case scenarios
- **Example:** Run 5 pods most of the month, scale to 30 only on salary day
- Significant cloud cost savings (possibly 80% reduction in idle capacity costs)

### 4. Agility for Growth
- Easily absorb sudden demand from new product launches or market expansion
- **Example:** If FinPay partners with a major bank overnight, auto-scaling cushions the surge of new users
- No emergency late-night deployments needed

## Real-World Flow (FinPay Example)

### Timeline: Salary Day Midnight

**23:50 - Pre-surge**
- Pods: 3
- Traffic: 8,000 transactions/hour
- CPU: 35%
- Status: Normal operations

**00:00 - Surge begins**
- Pods: 3
- Traffic: 50,000 transactions/hour (sudden spike)
- CPU: 95%
- HPA: Detects high CPU, starts scaling

**00:02 - First scale-up**
- Pods: 6 (doubled)
- Traffic: 100,000 transactions/hour
- CPU: 88%
- HPA: Still above threshold, continues scaling

**00:05 - Peak scaling**
- Pods: 20
- Traffic: 200,000 transactions/hour
- CPU: 72%
- HPA: Below threshold, stops scaling

**00:15 - Stabilized**
- Pods: 20
- Traffic: 180,000 transactions/hour
- CPU: 68%
- Status: Handling peak load smoothly

**14:00 - Post-peak scale-down**
- Pods: 5
- Traffic: 12,000 transactions/hour
- CPU: 42%
- Status: Cost-optimized for normal load

## Troubleshooting

**HPA shows `<unknown>` for targets:**
- Metrics server not installed or not running
- Resource requests not set on deployment
- Solution: Verify metrics-server and add resource requests

**Pods not scaling up despite high CPU:**
- Check if max replicas reached
- Verify HPA is watching correct deployment
- Check for resource constraints on nodes

**Pods scaling too aggressively:**
- Adjust `averageUtilization` threshold (try 85% instead of 80%)
- Add `behavior` section to HPA for slower scale-up

**Expensive: Pods not scaling down:**
- HPA has a default 5-minute cooldown for scale-down
- This is intentional to prevent flapping
- Can be adjusted in HPA `behavior` section

## Advanced: Custom Metrics

HPA can also scale based on custom metrics like:
- Request rate (requests per second)
- Queue depth (messages waiting to be processed)
- Response latency (99th percentile response time)

Example: Scale based on requests per second:
```yaml
metrics:
- type: Pods
  pods:
    metric:
      name: http_requests_per_second
    target:
      type: AverageValue
      averageValue: "1000"
```

## Key Takeaway

Auto-scaling is essential for cloud-native fintech applications. It ensures **reliability during unpredictable traffic spikes** while maintaining **cost efficiency during quiet periods**.

FinPay Wallet can handle salary day surges without manual intervention and without paying for idle capacity the rest of the month. This combination of performance and cost control is what makes cloud-native architecture powerful for fintech businesses.