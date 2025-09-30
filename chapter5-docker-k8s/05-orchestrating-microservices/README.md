# Orchestrating Microservices with Kubernetes

This example demonstrates orchestrating multiple microservices in Kubernetes.

## Services

1. **Payment Service** - 3 replicas on port 8080
2. **Fraud Detection Service** - 2 replicas on port 9090
3. **Notification Service** - 2 replicas on port 7070

## Architecture

```
User Request → Payment Service → Fraud Service
                              ↓
                       Notification Service
```

## Deployment

### 1. Create namespace
```bash
kubectl create namespace finpay
kubectl config set-context --current --namespace=finpay
```

### 2. Deploy all services
```bash
kubectl apply -f fraud.yaml
kubectl apply -f notification.yaml
kubectl apply -f payment.yaml
```

### 3. Verify deployment
```bash
kubectl get all
kubectl get deployments
kubectl get pods
kubectl get services
```

## Service Discovery

Services can communicate using DNS names:
- Payment → Fraud: `http://fraud-service:9090`
- Payment → Notification: `http://notification-service:7070`

## Testing

### Port forward to payment service
```bash
kubectl port-forward svc/payment-service 8080:8080
```

### Test
```bash
curl http://localhost:8080/pay
```

## Scaling

### Scale individual services
```bash
# Scale payment service
kubectl scale deployment payment-deployment --replicas=5

# Scale fraud service
kubectl scale deployment fraud-deployment --replicas=4

# Check status
kubectl get pods
```

## Monitoring

### View logs
```bash
# Payment service logs
kubectl logs -l app=payment --tail=50 -f

# Fraud service logs
kubectl logs -l app=fraud --tail=50 -f
```

### Check service endpoints
```bash
kubectl describe svc payment-service
kubectl describe svc fraud-service
```

## Clean Up

```bash
kubectl delete -f payment.yaml
kubectl delete -f fraud.yaml
kubectl delete -f notification.yaml
kubectl delete namespace finpay
```

## Options to Run Locally

### Option A: Minikube
```bash
minikube start --driver=docker
kubectl apply -f .
minikube service payment-service --namespace finpay
```

### Option B: kind
```bash
kind create cluster --name finpay
kubectl apply -f .
kubectl port-forward svc/payment-service 8080:8080
```

### Option C: Docker Desktop Kubernetes
1. Enable Kubernetes in Docker Desktop settings
2. Apply manifests: `kubectl apply -f .`
3. Port forward: `kubectl port-forward svc/payment-service 8080:8080`