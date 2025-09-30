# Cloud-Native Design Patterns with Go (Fintech Examples)

This repository contains Go code examples and practice exercises used in the book/course:
"Cloud-Native Design Patterns with Go for Fintech Systems".

The code illustrates cloud-native architecture principles (e.g., microservices, sidecar pattern, event-driven systems, service discovery, loose coupling) with practical Go snippets so that learners — especially beginners — can connect theory with real-world practice.

## 📂 Repository Structure

```
.
├── chapter3-go-fundamentals/     # Go basics (syntax, structs, interfaces, concurrency)
├── chapter4-patterns/            # Cloud-native design patterns (sidecar, API gateway, etc.)
├── chapter5-docker-k8s/          # Dockerization and Kubernetes orchestration
├── chapter6-principles/          # Cloud-native design principles (scalability, resilience)
├── chapter7-scalability/         # Horizontal/vertical scaling, auto-scaling
├── chapter8-loose-coupling/      # Event-driven, service discovery, API contracts
└── README.md                     # This file
```

Each chapter folder contains small, focused Go programs (not full apps), so learners can run and understand them step by step.

## 🚀 Running the Go Examples

**Install Go (v1.20+ recommended)**
[Download here](https://go.dev/dl/).

**Clone the repo**

```bash
git clone https://github.com/<your-org>/<repo-name>.git
cd <repo-name>
```

**Run any example**

```bash
cd chapter3-go-fundamentals/hello-world
go run main.go
```

**Build a binary**

```bash
go build -o app main.go
./app
```

## 🐳 Running with Docker

Some examples show how Go services can be containerized.

Typical workflow:

```bash
# Build image
docker build -t myapp .

# Run container
docker run -p 8080:8080 myapp

# Stop container
docker ps
docker stop <container-id>
```

## ☸️ Running with Kubernetes

Examples in `chapter5-docker-k8s/` provide Kubernetes YAML files.

Steps:

```bash
# Apply deployment
kubectl apply -f deployment.yaml

# Check pods
kubectl get pods

# Expose service
kubectl port-forward svc/myapp 8080:8080
```

## 🎯 Learning Objectives

By using this repository, learners will:

- Understand Go fundamentals (syntax, structs, interfaces, concurrency).
- Apply cloud-native design patterns like sidecar, API gateway, and event-driven systems.
- Practice Dockerizing and deploying Go services to Kubernetes.
- Learn scalability, resilience, and loose coupling through fintech-inspired case studies.
- Gain confidence in building real-world, cloud-native microservices in Go.

## 📝 Notes

- All code examples are kept simple and beginner-friendly.
- Each example includes comments explaining the logic step by step.
- Real-world fintech contexts (e.g., payments, fraud checks, notifications) are used so learners see practical relevance.