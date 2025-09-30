# 9. Continuous Delivery and Rapid Iteration

## Concepts

**Continuous Delivery (CD)** enables teams to:
- Deploy changes rapidly and safely
- Use feature flags to enable/disable features without redeployment
- Perform rolling updates with zero downtime
- Automate the entire pipeline from code commit to production

In FinTech, this is critical for:
- Launching new payment features (e.g., split bill) gradually
- Responding to market changes quickly
- Maintaining high availability during deployments

## Prerequisites

- Go 1.20+
- Docker
- Kubernetes cluster (minikube, kind, or Docker Desktop)
- kubectl
- GitHub account (for CI/CD)
- Docker Hub account

## How to Run

### Step 0: Local Development

```bash
cd chapter6-principles/09-continuous-delivery

# Run tests
go mod init payment
go test -v

# Run locally
go run main.go

# Test payment endpoint
curl -X POST http://localhost:8080/pay -d '{"amount": 100, "user": "alice"}'

# Test health endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

### Step 1: Build and Push Docker Image

```bash
# Build image
docker build -t <YOUR_DOCKERHUB_USERNAME>/payment:v1 .

# Test container locally
docker run -p 8080:8080 -e FEATURE_SPLIT_BILL=true <YOUR_DOCKERHUB_USERNAME>/payment:v1

# Push to Docker Hub
docker login
docker push <YOUR_DOCKERHUB_USERNAME>/payment:v1
```

### Step 2: Deploy to Kubernetes

Start your cluster:

```bash
# For minikube
minikube start

# For kind
kind create cluster

# For Docker Desktop
# Enable Kubernetes in Docker Desktop settings
```

Deploy the application:

```bash
# Edit k8s/deployment.yaml: replace <YOUR_DOCKERHUB_USERNAME> and IMAGE_TAG with v1
kubectl apply -f k8s/deployment.yaml

# Check deployment
kubectl get deployments
kubectl get pods
kubectl get services

# Access the service
# For minikube
minikube service payment --url

# For kind or Docker Desktop
kubectl port-forward svc/payment 8080:80
```

### Step 3: Test Deployment

```bash
# Test payment
curl -X POST http://localhost:8080/pay -d '{"amount": 100, "user": "alice"}'

# Expected response (feature flag is false by default):
# {"status":"success","request":{"amount":100,"feature":"split_bill_disabled","user":"alice"}}
```

### Step 4: Rolling Update with Feature Flag

Enable the split bill feature:

```bash
# Edit k8s/deployment.yaml: change FEATURE_SPLIT_BILL from "false" to "true"
# Then apply:
kubectl apply -f k8s/deployment.yaml

# Watch rolling update
kubectl rollout status deployment/payment

# Test again
curl -X POST http://localhost:8080/pay -d '{"amount": 100, "user": "alice"}'

# Now response shows:
# {"status":"success","request":{"amount":100,"feature":"split_bill_enabled","user":"alice"}}
```

### Step 5: Rollback (if needed)

```bash
# Rollback to previous version
kubectl rollout undo deployment/payment

# Check rollout history
kubectl rollout history deployment/payment
```

## Setting Up CI/CD with GitHub Actions

### Step 6: Prepare GitHub Repository

```bash
# Create a new GitHub repository
# Push your code:
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/<YOUR_USERNAME>/<REPO_NAME>.git
git push -u origin main
```

### Step 7: Configure GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions → New repository secret

Add these secrets:
1. `DOCKERHUB_USERNAME` - Your Docker Hub username
2. `DOCKERHUB_TOKEN` - Docker Hub access token (create at hub.docker.com → Account Settings → Security)
3. `KUBECONFIG` - Base64-encoded kubeconfig file:
   ```bash
   cat ~/.kube/config | base64
   ```

### Step 8: Add GitHub Actions Workflow

The workflow file `.github/workflows/cicd.yaml` is already in place. It will:
1. Checkout code
2. Set up Go
3. Run tests
4. Build Docker image (tagged with commit SHA)
5. Push to Docker Hub
6. Deploy to Kubernetes

### Step 9: Trigger the Pipeline

```bash
# Make a change to main.go (e.g., add a log message)
git add .
git commit -m "Update payment service"
git push

# Go to GitHub repository → Actions tab to watch the pipeline
```

## Real-World Flow (FinPay Example)

### Scenario: Launch Split Bill Feature Safely

1. **Developer commits code** with feature flag `FEATURE_SPLIT_BILL=false`
2. **CI/CD pipeline**:
   - Runs unit tests
   - Builds Docker image with commit SHA as tag
   - Deploys to Kubernetes with rolling update strategy
3. **Zero downtime**: Old pods serve traffic while new pods start
4. **Gradual rollout**:
   - Week 1: Deploy with flag=false (feature hidden)
   - Week 2: Enable for 10% users (canary deployment)
   - Week 3: Enable for all users (flag=true)
5. **Rollback**: If issues arise, revert deployment or toggle flag off

## Benefits

- **Faster time to market**: Deploy multiple times per day
- **Reduced risk**: Rolling updates and feature flags minimize blast radius
- **Quick recovery**: Instant rollback capability
- **Automation**: No manual steps, consistent process
- **Confidence**: Automated tests catch issues early

## Troubleshooting

**Pipeline fails at "Run tests" step**:
- Check test logic in `main_test.go`
- Ensure Go module is initialized correctly

**Image push fails**:
- Verify Docker Hub credentials in GitHub secrets
- Check Docker Hub rate limits

**Deployment fails**:
- Ensure `KUBECONFIG` secret is valid and base64-encoded
- Check Kubernetes cluster is accessible
- Verify image tag matches in deployment.yaml

**Service not accessible**:
- Check service type (LoadBalancer, NodePort, or port-forward)
- For minikube: Use `minikube service payment --url`
- For kind: Use `kubectl port-forward svc/payment 8080:80`

## Key Takeaway

Continuous Delivery is not just about automation—it's about enabling rapid, safe iteration. Feature flags and rolling updates give you control over when and how features reach users, which is essential in fintech where stability and compliance matter.