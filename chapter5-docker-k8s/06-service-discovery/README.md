# Service Discovery and Networking in Kubernetes

This example demonstrates service discovery using DNS in Kubernetes.

## Concept

Services communicate using stable DNS names instead of Pod IPs:
- Payment Service calls: `http://fraud-service:9090`
- DNS name resolves to current healthy Pods automatically

## How It Works

1. **CoreDNS**: Built-in DNS server resolves service names
2. **Service**: Provides stable VIP + DNS name for groups of Pods
3. **Cluster Networking**: All Pods can talk to all other Pods by default

## Deployment Steps

### 1. Create namespace
```bash
kubectl create namespace finpay
kubectl config set-context --current --namespace=finpay
```

### 2. Deploy Fraud service first
```bash
kubectl apply -f fraud.yaml
kubectl get deploy,po,svc -l app=fraud
```

### 3. Deploy Payment service (talks to Fraud)
```bash
kubectl apply -f payment.yaml
kubectl get deploy,po,svc -l app=payment
```

## Test Service Discovery

### Option A: From inside the cluster (debug pod)
```bash
# Start a temporary shell
kubectl run -it net-debug --image=busybox --restart=Never -- sh

# Inside the shell, call Fraud by DNS
wget -qO- http://fraud-service:9090/health

# Test FQDN (Fully Qualified Domain Name)
wget -qO- http://fraud-service.finpay.svc.cluster.local:9090/health

# Exit
exit
```

### Option B: Through Payment service
```bash
# Port forward Payment
kubectl port-forward svc/payment-service 8080:8080

# Test (in another terminal)
curl http://localhost:8080/pay

# View Payment logs to see it calling Fraud
POD=$(kubectl get pods -l app=payment -o jsonpath='{.items[0].metadata.name}')
kubectl logs -f "$POD"
```

## Verify DNS Resolution

```bash
# Get service details
kubectl get svc fraud-service -o wide
kubectl describe svc fraud-service

# Check endpoints (shows actual Pod IPs)
kubectl get endpoints fraud-service
```

## Scale and Test Stability

### Scale Fraud up
```bash
kubectl scale deployment fraud-deployment --replicas=4
kubectl get pods -l app=fraud -w
```

### Test - traffic still routes correctly
```bash
curl http://localhost:8080/pay
```

### Delete a Fraud pod
```bash
FRAUDPOD=$(kubectl get pods -l app=fraud -o jsonpath='{.items[0].metadata.name}')
kubectl delete pod "$FRAUDPOD"
kubectl get pods -l app=fraud -w
```

Service continues routing to healthy pods automatically.

## DNS Names Format

Short name (same namespace):
```
fraud-service:9090
```

FQDN (Fully Qualified Domain Name):
```
fraud-service.finpay.svc.cluster.local:9090
```

## Clean Up

```bash
kubectl delete -f payment.yaml
kubectl delete -f fraud.yaml
kubectl delete pod net-debug --ignore-not-found
kubectl delete ns finpay
```

## Key Concepts

- **Service DNS**: Stable name even when Pods restart
- **Load Balancing**: Automatically distributes traffic
- **Service Discovery**: No hardcoded IPs needed
- **CoreDNS**: Kubernetes DNS server
- **ClusterIP**: Internal-only service type