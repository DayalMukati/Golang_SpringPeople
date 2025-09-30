# CI/CD Integration with Docker and Kubernetes

This example demonstrates a complete CI/CD pipeline using GitHub Actions.

## What is CI/CD?

### CI (Continuous Integration)
- Every code push triggers automated build and tests
- Catches bugs early
- Ensures code quality

### CD (Continuous Delivery/Deployment)
- Automated packaging into Docker images
- Automated deployment to Kubernetes
- Frequent, safe releases

## Files

- **main.go** - Go payment service
- **Dockerfile** - Multi-stage Docker build
- **k8s/payment-deployment.yaml** - Kubernetes manifests
- **.github/workflows/deploy.yaml** - GitHub Actions pipeline

## Pipeline Flow

1. Developer pushes code to GitHub
2. GitHub Actions triggers:
   - Run tests (`go test`)
   - Build Docker image
   - Push to Docker Hub
   - Deploy to Kubernetes
3. Kubernetes performs rolling update (zero downtime)

## Setup Steps

### 1. Create GitHub Repository

```bash
# Create new repo: finpay-payment-service
# Push files to GitHub
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/<username>/finpay-payment-service.git
git push -u origin main
```

### 2. Add GitHub Secrets

Go to: Settings → Secrets → Actions → New repository secret

Add:
- **DOCKER_USERNAME**: Your Docker Hub username
- **DOCKER_PASSWORD**: Your Docker Hub token/password
- **KUBE_CONFIG** (optional): For remote cluster deployment

### 3. Deploy to Local Cluster

#### Using Minikube
```bash
minikube start --driver=docker
kubectl apply -f k8s/payment-deployment.yaml
```

#### Using kind
```bash
kind create cluster --name finpay
kubectl apply -f k8s/payment-deployment.yaml
```

### 4. Test Deployment

```bash
# Port forward
kubectl port-forward svc/payment-service 8080:8080

# Test
curl http://localhost:8080/pay
```

### 5. Trigger CI/CD

#### Make a code change
```bash
# Edit main.go - change message
vim main.go
# Change: "Payment processed successfully!" → "Payment completed!"

# Commit and push
git add main.go
git commit -m "Update payment message"
git push origin main
```

#### Watch GitHub Actions
- Go to your repo → Actions tab
- See the pipeline running
- Steps: Checkout → Build → Test → Push → Deploy

#### Verify deployment
```bash
# Check rollout status
kubectl rollout status deployment/payment-deployment

# Check pods (new pods with updated image)
kubectl get pods -l app=payment

# Test updated service
curl http://localhost:8080/pay
# Should show: "Payment completed!"
```

## GitHub Actions Workflow Explained

```yaml
name: CI/CD for Payment Service

on:
  push:
    branches: ["main"]  # Trigger on push to main

jobs:
  build-test-deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code           # Get source code
    - name: Set up Go               # Install Go
    - name: Run tests               # go test ./...
    - name: Build Docker image      # docker build
    - name: Push Docker image       # docker push
    - name: Deploy to Kubernetes    # kubectl apply
```

## Rolling Update

Kubernetes performs zero-downtime updates:
1. Creates new Pods with updated image
2. Waits for new Pods to be ready
3. Terminates old Pods gradually
4. Service continues running throughout

## Monitor Deployment

```bash
# Watch rollout
kubectl rollout status deployment/payment-deployment

# View history
kubectl rollout history deployment/payment-deployment

# Rollback if needed
kubectl rollout undo deployment/payment-deployment
```

## Clean Up

```bash
kubectl delete -f k8s/payment-deployment.yaml
minikube delete  # or: kind delete cluster --name finpay
```

## Real-World Benefits

### Without CI/CD
- Manual builds
- Manual testing
- Error-prone deployments
- Infrequent releases (monthly)

### With CI/CD
- Automated builds
- Automated tests
- Safe deployments
- Multiple releases per day
- Fast bug fixes (minutes, not days)

## Example Fintech Scenario

A developer fixes a rounding bug in payment calculations:
1. Push code to GitHub (1 minute)
2. Pipeline runs tests, builds, deploys (3-5 minutes)
3. Bug fix is live in production
4. Total time: ~5 minutes vs. hours/days manually