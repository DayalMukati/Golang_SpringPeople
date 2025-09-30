# Kubernetes Fundamentals (Pods, Services, Deployments)

This folder contains basic Kubernetes manifests demonstrating core concepts.

## Files

1. **pod.yaml** - Single Pod definition
2. **service.yaml** - Service for load balancing
3. **deployment.yaml** - Deployment managing multiple Pods

## Concepts

### Pod
- Smallest deployable unit in Kubernetes
- Wraps one or more containers
- Has unique IP address
- Ephemeral (can be created/destroyed)

### Service
- Provides stable endpoint for Pods
- Load balances traffic across Pods
- DNS name remains constant even if Pods restart
- Types: ClusterIP, NodePort, LoadBalancer

### Deployment
- Manages a set of Pods
- Defines desired state (replicas, image, etc.)
- Ensures desired state is maintained
- Supports rolling updates

## How to Use

### Prerequisites
```bash
# Install minikube or kind for local Kubernetes
# minikube: brew install minikube
# kind: brew install kind

# Start cluster
minikube start
# OR
kind create cluster
```

### Apply Manifests

```bash
# Create a Pod
kubectl apply -f pod.yaml
kubectl get pods

# Create a Service
kubectl apply -f service.yaml
kubectl get services

# Create a Deployment
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get pods
```

### Test the Deployment

```bash
# Scale the deployment
kubectl scale deployment payment-deployment --replicas=5

# Check status
kubectl get pods -l app=payment

# Port forward to access service
kubectl port-forward svc/payment-service 8080:8080

# Test (in another terminal)
curl http://localhost:8080/pay
```

### Clean Up

```bash
kubectl delete -f deployment.yaml
kubectl delete -f service.yaml
kubectl delete -f pod.yaml
```

## Key Commands

```bash
# Get resources
kubectl get pods
kubectl get services
kubectl get deployments

# Describe resources
kubectl describe pod <pod-name>
kubectl describe service <service-name>

# View logs
kubectl logs <pod-name>

# Execute commands in pod
kubectl exec -it <pod-name> -- /bin/sh
```