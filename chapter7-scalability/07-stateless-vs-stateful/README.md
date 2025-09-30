# 7. Stateless vs. Stateful Services in Scaling

## Concepts

### What Are Stateless and Stateful Services?

**Stateless Services:**
- Do not store data about previous requests
- Every request is independent — the server doesn't "remember" anything once it's done
- **Example:** A payment authorization API where each request contains all the info needed (amount, account, auth token)

**Stateful Services:**
- Depend on stored data or session history to work
- **Example:** A user session service that remembers Ravi's login, wallet balance, or transaction history between requests

**In simple terms:**
- **Stateless** = short memory (easy to clone and scale)
- **Stateful** = long memory (harder to scale, must manage data carefully)

### Why This Matters in Scaling

**Stateless services scale easily:**
- You can add 10 more pods, and any pod can serve any user
- Load balancers don't need to care which pod gets the request

**Stateful services are harder:**
- Data must be kept consistent across pods
- If Ravi's session is stored only on Pod 3 and Pod 3 crashes, his session is lost unless state is externalized

**That's why cloud-native design pushes for stateless services where possible, and uses specialized strategies for stateful ones.**

## FinPay Wallet Examples

### Stateless Example: Payment API

**Scenario:** Ravi sends ₹5,000 to his landlord

**How it works:**
- Request carries all info needed: sender ID, receiver ID, amount, and auth token
- Any pod can process it, hit the DB, and return success
- No pod needs to remember previous requests

**Scaling:** Simple — add pods, and the load balancer spreads traffic

### Stateful Example: Notification Service

**Scenario:** FinPay keeps an in-memory queue of SMS messages per user

**Problem:**
- If Ravi's SMS queue lives only in Pod 2 and Pod 2 crashes, his "Payment Sent" alert never gets delivered

**Scaling:** Harder — the state (queue) must be stored in a persistent store like Kafka or Redis to survive crashes

## Prerequisites

- Go 1.20+
- Basic understanding of functions and variables
- Conceptual understanding of distributed systems

## How to Run

### Step 1: Run Stateless Example

```bash
cd chapter7-scalability/07-stateless-vs-stateful

# Run stateless example
go run stateless_example.go
```

**Expected output:**
```
User: Ravi, Result: 6000
User: Meena, Result: 4000
User: Ravi, Result: 6000 (consistent)
```

**Observation:**
- Each call is independent
- Same input always produces same output
- No shared state between calls
- **Any pod can handle any request**

### Step 2: Run Stateful Example

```bash
# Run stateful example
go run stateful_example.go
```

**Expected output:**
```
Initial balance: 1000
After adding 500: 1500
After adding 300: 1800

Note: This pattern doesn't scale well across multiple pods!
```

**Observation:**
- Function remembers previous calls via global variable
- Each call modifies shared state
- **Problem:** If this runs in multiple pods, keeping balances consistent becomes a headache

### Step 3: Understand the Difference

**Stateless (AddBalance):**
```go
func AddBalance(user string, amount int) int {
    return amount + 1000  // Pure function, no memory
}
```
- Can run on Pod 1, Pod 2, or Pod 3 — doesn't matter
- Easy to scale horizontally
- No coordination needed

**Stateful (UpdateBalance):**
```go
var balance = 1000  // Global state

func UpdateBalance(amount int) int {
    balance += amount  // Modifies state
    return balance
}
```
- If Pod 1 has balance=1500 and Pod 2 has balance=1800, which is correct?
- Requires synchronization across pods
- Hard to scale horizontally

## Benefits of Stateless Services

### 1. Easy to Scale
- Stateless services don't depend on memory or local state — each request is self-contained
- If traffic spikes (like salary day when millions of payment requests flood in), new pods can be added instantly
- No need to worry about syncing data between pods, because each request carries everything it needs

**Fintech Example:**
- The Payment API in FinPay can scale from 3 to 30 pods at midnight
- Every pod processes transactions independently without any coordination

### 2. Fault Tolerance
- If a pod crashes, another pod can handle new requests immediately because there's no lost context
- Failures are invisible to users, since no critical information was stored inside the failed pod

**Fintech Example:**
- If Pod 7 dies in the middle of processing, Ravi's next payment request simply goes to Pod 8
- Pod 8 starts fresh and completes the transfer successfully

### 3. Simpler Load Balancing
- Since any pod can handle any request, load balancers can use simple algorithms like round robin
- There's no need for sticky sessions (where the same user must always return to the same pod)

**Fintech Example:**
- When Ravi and Meena check their balances, requests can go to Pod 1, Pod 2, or Pod 3 interchangeably
- All pods work the same way, so scaling and traffic routing are straightforward

## Challenges with Stateful Services

### 1. Consistency Problems
- When state is stored locally in a pod, keeping data consistent across pods is complicated
- **Example:** If Ravi's session lives in Pod 3 and Meena's in Pod 4, but a cross-user fraud check needs both, coordinating state becomes painful

### 2. Harder to Scale
- Scaling a stateful service isn't just about adding pods — you also need to decide which pod stores which data
- This usually requires sharding or partitioning, which adds complexity

**Fintech Example:**
- If session data is stored locally, FinPay must ensure Ravi always lands on Pod 3
- Adding new pods would require reshuffling who holds which users' sessions

### 3. Risk of Data Loss
- If a pod with in-memory state crashes, all of that data disappears
- For critical systems like payments, this risk is unacceptable

**Fintech Example:**
- If Ravi's "payment pending" queue is in Pod 2's memory and Pod 2 crashes, that payment may vanish without being processed
- **A nightmare for both user and regulator**

## Best Practice: Externalize State

Because of these risks, fintech apps externalize state to reliable systems:

### Where to Store State

**1. Databases (Postgres, MySQL)**
- For transactions, account balances, audit logs
- ACID guarantees ensure data integrity

**2. Redis**
- For caching session data, temporary locks, rate limits
- Fast in-memory access with optional persistence

**3. Kafka**
- For durable message queues, event streams
- Ensures messages aren't lost even if pods crash

**4. Object Storage (S3, GCS)**
- For documents, receipts, compliance records
- Durable and cost-effective for large files

### Result

Core services remain **stateless** and easy to scale, while state lives in **fault-tolerant systems** designed to handle persistence and recovery.

## Architecture Patterns

### Pattern 1: Stateless API + External Database

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  Pod 1  │     │  Pod 2  │     │  Pod 3  │
│(Stateless)   │(Stateless)   │(Stateless)
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┴───────────────┘
                     │
                ┌────▼────┐
                │Database │
                │(Postgres)
                └─────────┘
```

**Benefits:**
- Any pod can handle any request
- Easy horizontal scaling
- Database handles consistency

### Pattern 2: Stateless API + Redis Cache

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  Pod 1  │     │  Pod 2  │     │  Pod 3  │
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┴───────────────┘
                     │
            ┌────────┴────────┐
            │                 │
       ┌────▼────┐      ┌────▼─────┐
       │  Redis  │      │ Database │
       │ (Cache) │      │(Postgres)│
       └─────────┘      └──────────┘
```

**Benefits:**
- Fast reads from Redis
- Reduced database load
- Pods remain stateless

### Pattern 3: Stateful Service (Kubernetes StatefulSet)

For services that **must** be stateful (e.g., databases themselves):

```
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  Pod-0       │   │  Pod-1       │   │  Pod-2       │
│  (Primary)   │──>│  (Replica)   │──>│  (Replica)   │
│              │   │              │   │              │
│  ┌────────┐  │   │  ┌────────┐  │   │  ┌────────┐  │
│  │ Volume │  │   │  │ Volume │  │   │  │ Volume │  │
│  └────────┘  │   │  └────────┘  │   │  └────────┘  │
└──────────────┘   └──────────────┘   └──────────────┘
```

**Use cases:**
- Database clusters
- Message queue brokers
- Distributed caches

## Real-World Decision Matrix

| Service Type | State Storage | Scaling Strategy | Example |
|-------------|---------------|------------------|---------|
| Payment API | Database | Stateless, horizontal | Payment processing |
| Balance Check | Redis + DB | Stateless, horizontal | Balance queries |
| Session Management | Redis | Stateless, horizontal | User sessions |
| Transaction History | Database | Stateless, horizontal | History queries |
| Fraud Detection | Redis (rules) + DB | Stateless, horizontal | Real-time fraud checks |
| Notification Queue | Kafka | Stateless workers | SMS/email sending |
| Database | Local volumes | StatefulSet | Postgres, MongoDB |

## Key Takeaway

**Design principle:** Make services stateless by default. When state is unavoidable, externalize it to purpose-built systems (databases, caches, queues).

For FinPay Wallet:
- **Payment API:** Stateless — scales from 3 to 30 pods seamlessly
- **Session data:** Stored in Redis — any pod can access
- **Transaction records:** Stored in Postgres — ACID guarantees
- **Notification queue:** Stored in Kafka — survives pod crashes

This architecture ensures:
- **Easy scaling** during salary day surges
- **Fault tolerance** when pods crash
- **Data integrity** for financial transactions
- **Operational simplicity** for engineering teams

**Remember:** Stateless services are the foundation of cloud-native scalability. Externalize state, and your system becomes resilient, scalable, and maintainable.