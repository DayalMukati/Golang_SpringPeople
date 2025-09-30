# 5. Event-Driven Architecture for Decoupling Services

## Concepts

### What Is Event-Driven Architecture?

**Event-driven architecture (EDA)** is a design style where services don't talk to each other directly. Instead, they communicate by emitting and reacting to **events**.

**Event** → A fact about something that happened in the system

**Examples in fintech:**
- `"TransactionCreated"` → when Ravi transfers ₹1000 to his landlord
- `"FraudCheckFailed"` → when a risky transaction is detected
- `"NotificationSent"` → when an SMS or email alert goes out

**Publisher** → The service that generates the event and pushes it into an event bus

**Example:** Payment Service publishes a `TransactionCreated` event after a payment request

**Subscriber** → Any service that listens for that event and reacts accordingly

**Example:** Notification Service subscribes to `TransactionCreated` and sends an SMS to Ravi

### Key Difference from Direct Communication

**Instead of:** "Hey Fraud, here's a payment — check it"

**EDA says:** "A payment was created" (whoever cares can listen and react)

### Event Bus / Streaming Platforms

To carry these events, systems use:
- **Apache Kafka**
- **RabbitMQ**
- **AWS SNS/SQS**
- **Google Pub/Sub**

Think of them as the **postal system** of cloud applications:
- Publishers drop letters (events) into the mailbox (bus)
- Subscribers open the mailbox when they're ready and read the letters

## Why Synchronous Communication Falls Short

In older or tightly coupled systems, services use **synchronous communication**: one service directly calls another's API and waits for a response.

### Problems with Synchronous Communication

**1. Tight Coupling**
- Payment must know where Fraud lives: its URL, network address, protocol
- If Fraud moves to another server or changes its API version, Payment breaks
- **Example:** FinPay upgrades Fraud to a new cluster. Suddenly, Payments start failing because URLs are outdated

**2. Blocking**
- If Fraud takes 2-3 seconds to check a transaction, Payment must wait that long before responding to the user
- Worse, if Fraud is down, Payments cannot be processed at all
- **Example:** On salary day, Ravi's payment spins for 5 seconds because Fraud is slow, frustrating him

**3. Cascading Failures**
- In synchronous chains, one failure spreads like dominoes
- If Notification Service crashes, Fraud may be fine, but Payment flow still fails because it's all linked together
- **Example:** Notification outage leads to failed payments, even though money transfers should work independently

**4. Scaling Issues**
- Synchronous systems force all services to scale together, even if only one needs more capacity
- **Example:** Payment spikes on salary day. With synchronous coupling, Fraud and Notification must also scale unnecessarily, wasting resources

**This is why synchronous communication isn't cloud-native.** In systems that must handle traffic spikes, partial failures, and frequent updates, synchronous calls become bottlenecks.

Event-driven communication fixes this by letting services **publish and react independently** — scaling, failing, and evolving at their own pace.

## How Event-Driven Architecture Works (Step by Step)

### 1. Something Happens
Ravi makes a rent payment

### 2. Publisher Emits an Event
Payment Service publishes:
```json
{
  "event": "TransactionCreated",
  "transactionId": "TXN123",
  "amount": 1000
}
```

### 3. Event Bus Delivers the Event
The event bus (e.g., Kafka) acts like a postal service

### 4. Subscribers React
- **Fraud Service** checks for fraud
- **Notification Service** sends SMS/Email
- **Reporting Service** logs the transaction

**Key:** The publisher (Payment) doesn't need to know who listens — it only says "Here's what happened."

## FinPay Wallet Example

In FinPay Wallet:

**Payment Service** publishes `TransactionCreated`

**Subscribers:**
- **Fraud Service** listens and checks if transaction is suspicious
- **Notification Service** listens and sends alerts
- **Reporting Service** listens and logs it for monthly statements

**Tomorrow:**
- A **Loyalty Service** could subscribe to same event to give cashback
- **Without touching Payment at all**

The system evolves by **adding subscribers**, not by modifying the publisher.

## Prerequisites

- Go 1.20+
- Understanding of channels (Go's concurrency primitive)
- Understanding of publish-subscribe pattern

## How to Run

### Step 1: Run the Event Demo

```bash
cd chapter8-loose-coupling/05-event-driven

# Run the demo
go run event_demo.go
```

**Expected output:**
```
Payment: publishing TransactionCreated
Fraud: checking transaction TXN123
```

**Observation:**
- Payment publishes an event and moves on
- Fraud subscribes and reacts when event arrives
- Payment doesn't know Fraud exists
- Fraud doesn't know who published the event

### Step 2: Understand the Code

**Event Structure:**
```go
type Event struct {
    Name string
    Data string
}
```

Defines the **contract** — what events look like.

**Publisher (Payment Service):**
```go
func payment(events chan Event) {
    fmt.Println("Payment: publishing TransactionCreated")
    events <- Event{Name: "TransactionCreated", Data: "TXN123"}
}
```

- Creates an event
- Publishes it to the event bus (channel)
- **Doesn't wait, doesn't know who's listening**

**Subscriber (Fraud Service):**
```go
func fraud(events chan Event) {
    for e := range events {
        if e.Name == "TransactionCreated" {
            fmt.Println("Fraud: checking transaction", e.Data)
        }
    }
}
```

- Listens to the event bus
- Filters for events it cares about (`TransactionCreated`)
- Reacts independently

**Main Function:**
```go
func main() {
    events := make(chan Event, 1)

    go fraud(events)  // Subscriber runs in background
    payment(events)   // Publisher sends event

    select {}  // Keep running
}
```

- Creates event bus (channel)
- Starts Fraud as subscriber
- Payment publishes event
- System runs independently

**In real fintech systems, Kafka/RabbitMQ replaces this Go channel.**

## Synchronous vs Event-Driven Comparison

### Synchronous (Tight Coupling)

```
Payment → Fraud → Notification → Reporting
  (each waits for previous)
```

**Problems:**
- Payment must know Fraud's URL
- If Fraud slow → Payment slow
- If Notification down → Payment fails
- Hard to add new consumers

### Event-Driven (Loose Coupling)

```
Payment → Event Bus → Fraud
                   → Notification
                   → Reporting
                   → Loyalty (new)
```

**Benefits:**
- Payment doesn't know who's listening
- Services react independently
- Can add new subscribers without touching Payment
- Services can fail independently

## Real-World Flow (FinPay Example)

### Scenario: Ravi Pays Rent

**Traditional Synchronous Flow:**

1. Payment calls Fraud API (waits 2s)
2. Payment calls Notification API (waits 1s)
3. Payment calls Reporting API (waits 0.5s)
4. Payment responds to Ravi: "Success" (after 3.5s total)

**Problems:**
- Ravi waits 3.5 seconds
- If any service fails, payment fails
- Payment knows about all services

**Event-Driven Flow:**

1. Payment publishes `TransactionCreated` event (0.1s)
2. Payment responds to Ravi: "Success" (0.1s total)
3. **In background:**
   - Fraud picks up event, checks transaction (2s)
   - Notification picks up event, sends SMS (1s)
   - Reporting picks up event, logs transaction (0.5s)

**Benefits:**
- ✅ Ravi gets instant response (0.1s vs 3.5s)
- ✅ If Notification fails, payment still succeeds
- ✅ Payment doesn't know about Fraud/Notification/Reporting
- ✅ Easy to add Loyalty service without touching Payment

### Adding New Feature (Loyalty Program)

**Synchronous approach:**
- Must update Payment Service code
- Add `callLoyaltyAPI()` function
- Deploy Payment Service
- Risk breaking existing flow

**Event-driven approach:**
- Create new Loyalty Service
- Subscribe to `TransactionCreated` event
- Deploy Loyalty Service independently
- **Zero changes to Payment Service**

**Result:** New features added without touching existing code = safer, faster

## Benefits of Event-Driven Architecture

### 1. Complete Decoupling
- Publisher doesn't know subscribers exist
- Subscribers don't know publisher exists
- Only contract is the event structure

### 2. Independent Evolution
- Add new subscribers anytime
- Remove old subscribers anytime
- Update subscriber logic without touching publisher

### 3. Better Resilience
- If Notification fails, Fraud and Reporting still work
- Failed services can retry when they recover
- No cascading failures

### 4. Scalability
- Payment scales to 30 pods
- Fraud scales to 5 pods
- Notification scales to 10 pods
- Each at its own pace, no coordination

### 5. Replay and Audit
- Events are facts: "TransactionCreated at 2024-01-15 14:30"
- Can be stored and replayed for debugging
- Provides audit trail for compliance

### 6. Flexibility
- Same event can trigger multiple reactions
- Easy to add new features (Loyalty, Analytics, Compliance checks)
- Business logic grows without touching core services

## Event-Driven Patterns

### 1. Publish-Subscribe
```
Payment → Event Bus → Multiple Subscribers
```
All subscribers get the same event

### 2. Event Sourcing
```
Store all events as source of truth
Current state = replay all events
```
Benefits: Complete audit trail, time travel debugging

### 3. CQRS (Command Query Responsibility Segregation)
```
Commands (writes) → Event Bus → Queries (reads)
Separate write and read models
```
Benefits: Optimize reads and writes independently

## Event Design Best Practices

### 1. Events are Facts
- Name events in past tense: `TransactionCreated`, not `CreateTransaction`
- Events describe what happened, not commands

### 2. Events Should Be Immutable
- Once published, never change
- Append new events instead

### 3. Include Enough Context
```json
{
  "event": "TransactionCreated",
  "transactionId": "TXN123",
  "userId": "Ravi123",
  "amount": 5000,
  "currency": "INR",
  "timestamp": "2024-01-15T14:30:00Z"
}
```

### 4. Version Events
```json
{
  "event": "TransactionCreated",
  "version": "v2",
  "data": { ... }
}
```

### 5. Keep Events Small
- Don't embed entire objects
- Include IDs, let subscribers fetch details if needed

## Event Bus Technologies

### Apache Kafka
- High-throughput streaming
- Persistent events (can replay)
- Good for event sourcing

### RabbitMQ
- Traditional message broker
- Flexible routing
- Good for microservices

### AWS SNS + SQS
- Managed pub-sub
- Scales automatically
- Good for AWS-native apps

### Google Pub/Sub
- Global scale
- At-least-once delivery
- Good for GCP-native apps

## Key Takeaway

**Event-driven architecture is the ultimate form of loose coupling.**

For FinPay Wallet:
- **Without events:** Payment must call Fraud, Notification, Reporting directly → tight coupling, slow, fragile
- **With events:** Payment publishes `TransactionCreated`, others react independently → loose coupling, fast, resilient

**Design principle:** "Tell everyone what happened, let them decide what to do about it."

Benefits:
- ✅ Services evolve independently
- ✅ Easy to add new features (just add subscribers)
- ✅ Resilient (failures don't cascade)
- ✅ Scalable (services scale independently)
- ✅ Better user experience (instant responses)
- ✅ Complete audit trail (all events stored)

**"In event-driven systems, the publisher doesn't care who's listening, and the subscribers don't care who's publishing. That's loose coupling at its finest."**