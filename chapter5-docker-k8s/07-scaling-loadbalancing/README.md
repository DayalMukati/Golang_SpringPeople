# Scaling and Load Balancing Microservices

This example demonstrates horizontal scaling and load balancing in Kubernetes.

## Concepts

### Scaling
- **Scale Out (Horizontal)**: Add more Pods
- **Scale In**: Reduce Pods when demand drops
- **Manual Scaling**: `kubectl scale`
- **Auto Scaling**: Horizontal Pod Autoscaler (HPA)

### Load Balancing
- Services automatically distribute requests across Pods
- Removes failed Pods from rotation
- No single point of failure

## Lab: Scaling and Load Balancing

### Prerequisites
```bash
# Enable metrics server for HPA (minikube)
minikube addons enable metrics-server
```

### 1. Create namespace
```bash
kubectl create namespace finpay
kubectl config set-context --current --namespace=finpay
```

### 2. Deploy Payment Service (2 replicas)
```bash
kubectl apply -f payment.yaml
kubectl get deploy,po,svc
```

### 3. Test the service
```bash
# Port forward
kubectl port-forward svc/payment-service 8080:8080

# Test (in another terminal)
curl http://localhost:8080/pay
```

### 4. Observe load balancing
```bash
# Make several requests
for i in {1..5}; do curl http://localhost:8080/pay; done

# Check logs on each Pod
kubectl get pods -l app=payment
POD1=<first-pod-name>
POD2=<second-pod-name>

kubectl logs $POD1
kubectl logs $POD2
```

Requests are distributed between Pods.

### 5. Manual Scaling

#### Scale out to 5 replicas
```bash
kubectl scale deployment payment-deployment --replicas=5
kubectl get pods -l app=payment
```

#### Test again
```bash
for i in {1..10}; do curl http://localhost:8080/pay; done
```

Requests now spread across 5 Pods.

#### Scale back down
```bash
kubectl scale deployment payment-deployment --replicas=2
kubectl get pods -l app=payment -w
```

### 6. Auto Scaling (HPA)

#### Apply HPA
```bash
kubectl apply -f payment-hpa.yaml
kubectl get hpa
```

HPA will:
- Monitor CPU usage
- Scale between 2-10 replicas
- Add Pods if CPU > 70%

#### Watch HPA in action
```bash
kubectl get hpa -w
```

#### Generate load (optional)
```bash
# In a separate terminal, generate load
while true; do curl http://localhost:8080/pay; done
```

Watch Pods scale up:
```bash
kubectl get pods -l app=payment -w
```

### 7. Self-Healing Demo

#### Delete a Pod
```bash
POD=$(kubectl get pods -l app=payment -o jsonpath='{.items[0].metadata.name}')
kubectl delete pod $POD
```

#### Watch automatic recreation
```bash
kubectl get pods -l app=payment -w
```

Kubernetes automatically recreates the Pod to maintain desired replica count.

## Key Commands

### Manual Scaling
```bash
# Scale to specific number
kubectl scale deployment payment-deployment --replicas=10

# Check status
kubectl get pods -l app=payment
```

### View HPA Status
```bash
kubectl get hpa
kubectl describe hpa payment-hpa
```

### View Service Endpoints
```bash
kubectl get endpoints payment-service
kubectl describe svc payment-service
```

## Clean Up

```bash
kubectl delete -f payment-hpa.yaml
kubectl delete -f payment.yaml
kubectl delete ns finpay
```

## Real-World Scenario

During a festival sale:
- Traffic spikes from 1,000 req/min to 10,000 req/min
- HPA scales Payment from 2 → 10 Pods automatically
- Load balancer distributes across all 10 Pods
- System stays fast and reliable
- After sale, scales back down to save resources