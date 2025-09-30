# 1. Introduction to Loose Coupling

## Concepts

### What is Loose Coupling?

**Loose coupling** is a design principle that ensures different parts of a system can work together without being tightly dependent on each other's inner details.

**In a tightly coupled system:**
- Components directly depend on each other's code, database structures, or workflows
- A change in one often forces changes in others

**In a loosely coupled system:**
- Components interact through well-defined contracts (like APIs or events)
- They care only about **what** the other service does, not **how** it does it

### Analogy: Two Companies

**Tight coupling:**
- Share same office, same employees, same bank account
- If one collapses, the other goes down

**Loose coupling:**
- Sign a clear contract
- Each manages its own operations
- If one upgrades its tools, the other doesn't need to change

## Why Is Loose Coupling Important in Cloud-Native Systems?

Cloud-native systems are built from many small services (microservices). If these services are tightly linked:

❌ **Problems:**
- Scaling becomes hard (all services must scale together)
- Changes become slow (every release requires coordination)
- Failures spread quickly (a bug in one service can crash the whole app)

✅ **Loose coupling solves these by ensuring:**
- Services can evolve independently
- Teams can deploy at their own pace
- Failures in one service don't stop the whole system

## FinPay Wallet Example

FinPay, a digital wallet app, with four services:
1. **Payment Service** → Handles money transfers
2. **Fraud Detection Service** → Flags suspicious activity
3. **Notification Service** → Sends SMS/email alerts
4. **Reporting Service** → Generates transaction history

### Case 1: Tightly Coupled Design

**Architecture:**
- Payment Service directly calls Fraud with custom logic
- Fraud then calls Notification once approved
- Notification writes to Reporting after sending SMS

**Problems:**
- If Fraud changes its function signature, Payment must be updated
- If Notification goes down, Payments fail too
- Scaling is inefficient: if Payments need 30 pods on salary day, Fraud and Notification must scale unnecessarily

**Flow:**
```
Payment → Fraud → Notification → Reporting
  (if any fails, entire chain breaks)
```

### Case 2: Loosely Coupled Design

**Architecture:**
- Payment Service publishes an event → "TransactionCreated"
- Fraud, Notification, and Reporting subscribe to that event
- Each service runs independently

**Benefits:**
- If Notification fails, Payments still succeed. SMS can be retried later
- Reporting can change its database schema without breaking Payments
- Fraud can adopt machine learning while Payment service code stays untouched

**Flow:**
```
Payment → Event Bus → Fraud
                   → Notification
                   → Reporting
  (each subscriber independent)
```

**Result:** Loose coupling means independence + resilience + speed

## Prerequisites

- Go 1.20+
- Basic understanding of interfaces
- Understanding of microservices architecture

## How to Run

### Step 1: Run Tightly Coupled Example

```bash
cd chapter8-loose-coupling/01-introduction

# Run tight coupling example
go run tight_coupling.go
```

**Observation:**
- `ProcessPayment` knows how fraud and notifications are done
- Any change in fraud or notification breaks this code
- Hard to test in isolation
- Hard to replace implementations

### Step 2: Run Loosely Coupled Example

```bash
# Run loose coupling example
go run loose_coupling.go
```

**Expected output:**
```
Notification: Payment processed successfully
```

**Observation:**
- `ProcessPayment` only depends on interfaces, not implementations
- Fraud and Notification can be swapped, replaced, or upgraded independently
- Easy to test with mock implementations
- Easy to extend with new implementations

### Step 3: Understand the Difference

**Tight Coupling:**
```go
func ProcessPayment(amount int) {
    fraudCheck(amount)      // Direct dependency
    sendNotification()      // Direct dependency
}
```

**Problems:**
- Hard-coded dependencies
- Cannot swap implementations
- Hard to test (must run real fraudCheck and sendNotification)
- Changes in one function break ProcessPayment

**Loose Coupling:**
```go
type FraudChecker interface {
    Check(amount int) bool
}

type Notifier interface {
    Send(message string)
}

func ProcessPayment(amount int, fc FraudChecker, n Notifier) {
    if fc.Check(amount) {
        n.Send("Payment processed successfully")
    }
}
```

**Benefits:**
- Depends on interfaces (contracts), not concrete implementations
- Can swap implementations easily (production vs test vs mock)
- Easy to test with mock implementations
- Changes in implementations don't break ProcessPayment

## Real-World Impact

### Tightly Coupled System

**Scenario:** Fraud team wants to upgrade from rule-based to ML-based detection

**Impact:**
- Must coordinate with Payment team
- Payment code needs updates
- Must deploy both services together
- Risk of breaking payments during upgrade
- Timeline: 2-3 weeks

### Loosely Coupled System

**Scenario:** Same upgrade (rule-based → ML-based)

**Impact:**
- Fraud team works independently
- Payment code unchanged (still uses FraudChecker interface)
- Deploy Fraud service independently
- Zero risk to payments
- Timeline: 3-5 days

## Key Principles

### 1. Depend on Abstractions, Not Concretions
```go
// ❌ Bad: Depends on concrete implementation
func ProcessPayment(amount int, checker SimpleFraudChecker) { ... }

// ✅ Good: Depends on interface
func ProcessPayment(amount int, checker FraudChecker) { ... }
```

### 2. Inject Dependencies
```go
// ❌ Bad: Creates dependencies internally
func ProcessPayment(amount int) {
    checker := SimpleFraudChecker{}  // hard-coded
    checker.Check(amount)
}

// ✅ Good: Dependencies injected from outside
func ProcessPayment(amount int, checker FraudChecker) {
    checker.Check(amount)
}
```

### 3. Program to Interfaces
```go
// ✅ Define behavior as interface
type FraudChecker interface {
    Check(amount int) bool
}

// ✅ Multiple implementations possible
type SimpleFraudChecker struct{}
type MLFraudChecker struct{}
type ThirdPartyFraudChecker struct{}
```

## Benefits Summary

| Aspect | Tight Coupling | Loose Coupling |
|--------|---------------|----------------|
| **Change Impact** | Ripples across system | Isolated to service |
| **Testing** | Hard (need real dependencies) | Easy (use mocks) |
| **Scaling** | All services together | Each service independently |
| **Deployment** | Coordinated releases | Independent releases |
| **Failure Impact** | Cascades to all services | Isolated to one service |
| **Technology Choice** | Must use same stack | Can use different stacks |
| **Team Autonomy** | Heavy coordination needed | Teams work independently |

## Key Takeaway

**Loose coupling is fundamental to cloud-native systems.**

For FinPay Wallet:
- Payment service doesn't need to know **how** Fraud detects suspicious activity
- It only needs to know **what** to send (amount, transaction ID) and **what** to expect back (fraudulent: true/false)
- This contract (interface) allows Fraud team to innovate (rules → ML → AI) without ever breaking Payment

**Design principle:** Program to interfaces, not implementations. This is the foundation of maintainable, scalable, cloud-native systems.