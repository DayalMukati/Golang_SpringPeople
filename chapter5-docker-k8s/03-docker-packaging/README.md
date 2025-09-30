# Building and Packaging Services with Docker

This example demonstrates how to build and containerize a Go microservice.

## Files

- **payment.go** - Simple payment service
- **Dockerfile** - Multi-stage Docker build

## Multi-Stage Build

The Dockerfile uses two stages:
1. **Builder stage** - Compiles Go code (needs Go installed)
2. **Runtime stage** - Smaller image with only the binary (uses debian:stable-slim)

## Docker Commands

### 1. Build the Image
```bash
docker build -t finpay/payment-service .
```

### 2. List Images
```bash
docker images
```

### 3. Run a Container
```bash
docker run -p 8080:8080 finpay/payment-service
```

### 4. Test the Service
```bash
curl http://localhost:8080/pay
# Output: Payment of 100 processed successfully!
```

### 5. Check Running Containers
```bash
docker ps
```

### 6. Stop a Container
```bash
docker stop <container-id>
```

### 7. Push to Docker Hub (optional)
```bash
docker tag finpay/payment-service <your-dockerhub-username>/payment-service:1.0
docker push <your-dockerhub-username>/payment-service:1.0
```

## Benefits of Multi-Stage Builds

- **Smaller images** - Final image only contains the binary, not Go toolchain
- **Faster deployments** - Less data to transfer
- **More secure** - Smaller attack surface
- **Separation of concerns** - Build environment separate from runtime