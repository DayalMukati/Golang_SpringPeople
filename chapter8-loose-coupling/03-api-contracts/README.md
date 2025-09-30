# 3. API Contracts and Service Boundaries

## Concepts

### What Are API Contracts?

An **API contract** defines how two services communicate. It specifies:
- The endpoint or message type (e.g., `POST /payments`)
- The input format (fields, types, required values)
- The output format (response, status codes, error messages)
- The rules of interaction (idempotency, authentication, rate limits)

**Think of it as a legal contract between services:**
- If you send me **this data** in the agreed format
- I promise to respond in **this way**

### What Are Service Boundaries?

A **service boundary** is the scope of responsibility for a service. Each service does one job well and exposes it through an API.

**Examples:**
- **Payments** → handle money transfers
- **Fraud** → detect risky transactions
- **Notifications** → send SMS/Email alerts
- **Reporting** → generate statements

**Boundaries prevent overlap:**
- Fraud shouldn't be sending emails
- Notifications shouldn't decide if a transaction is suspicious

### Why Are They Important?

❌ **Without clear contracts:**
- Services become tightly coupled
- Changes in one break others

❌ **Without boundaries:**
- Services grow messy
- Mix unrelated responsibilities

✅ **With both in place:**
- Systems stay stable, predictable, and easy to evolve

## Fintech Example: FinPay Wallet

Imagine Ravi sends money to his landlord.

### 1. Payment Service API Contract

**Endpoint:** `POST /payments`

**Input:**
```json
{
  "fromAccount": "Ravi123",
  "toAccount": "Landlord456",
  "amount": 1000
}
```

**Output:**
```json
{
  "status": "SUCCESS",
  "transactionId": "TXN789"
}
```

### 2. Fraud Service API Contract

**Endpoint:** `POST /fraud/check`

**Input:**
```json
{
  "transactionId": "TXN789",
  "amount": 1000
}
```

**Output:**
```json
{
  "fraudulent": false
}
```

### 3. Notification Service API Contract

**Endpoint:** `POST /notifications/send`

**Input:**
```json
{
  "to": "Ravi",
  "message": "Payment successful"
}
```

**Output:**
```json
{
  "status": "DELIVERED"
}
```

**Each service has:**
- One clear boundary
- An API that defines how others interact with it

**Key benefit:** If Fraud changes its internal logic, it doesn't affect Payment or Notification as long as the API contract remains the same.

## Prerequisites

- Go 1.20+
- Basic understanding of HTTP APIs
- Understanding of JSON serialization

## How to Run

### Step 1: Run the Payment API

```bash
cd chapter8-loose-coupling/03-api-contracts

# Run the payment API
go run payment_api.go
```

The server starts on port 8080.

### Step 2: Test the API Contract

**In another terminal:**

```bash
# Test the payment endpoint
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{
    "fromAccount": "Ravi123",
    "toAccount": "Landlord456",
    "amount": 5000
  }'
```

**Expected output:**
```json
{
  "status": "SUCCESS",
  "transactionId": "TXN12345"
}
```

### Step 3: Understand the Code

**Input Contract (PaymentRequest):**
```go
type PaymentRequest struct {
    FromAccount string  `json:"fromAccount"`
    ToAccount   string  `json:"toAccount"`
    Amount      float64 `json:"amount"`
}
```

Defines what Payment expects:
- `fromAccount` (string)
- `toAccount` (string)
- `amount` (float64)

**Output Contract (PaymentResponse):**
```go
type PaymentResponse struct {
    Status        string `json:"status"`
    TransactionID string `json:"transactionId"`
}
```

Defines what Payment guarantees back:
- `status` (string)
- `transactionId` (string)

**Handler enforces the contract:**
```go
func handlePayment(w http.ResponseWriter, r *http.Request) {
    var req PaymentRequest
    _ = json.NewDecoder(r.Body).Decode(&req)

    // Process payment
    res := PaymentResponse{
        Status:        "SUCCESS",
        TransactionID: "TXN12345",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}
```

**Result:** Other services can confidently call `/payments` knowing the request/response structure is stable.

## Benefits of API Contracts and Service Boundaries

### 1. Stability

An API contract acts like a **promise**. As long as the service keeps that promise (same request/response structure), other services don't care about what happens inside.

**Why it matters:**
- Services can evolve without breaking others
- Reduces fear of "breaking the system" every time a change is made

**Example:**
- Fraud Detection initially uses simple rules: block payments over ₹50,000 made at night
- Later, Fraud team upgrades to ML model with advanced risk scoring
- Since Fraud API still accepts same input (transactionId, amount) and returns fraudulent: true/false, Payment Service doesn't need to change

**Result:** Teams innovate inside their boundary without destabilizing the system

### 2. Clarity

Service boundaries ensure each service has a **single, well-defined responsibility**. No mixing of unrelated features.

**Why it matters:**
- Keeps systems easier to maintain
- Avoids overlapping responsibilities and messy ownership

**Example:**
- Payment Service = transfer money
- Fraud Service = check risk
- Notification Service = send alerts
- Reporting Service = generate statements

In a tight system, Payment might also check fraud, send SMS, and log reports — a maintenance nightmare. Clear contracts stop that.

**Result:** Each team owns one job, making system simpler and cleaner

### 3. Independent Development

With contracts in place, teams can develop and deploy **independently**. As long as the contract is respected, others don't need to know about internal changes.

**Why it matters:**
- Teams move faster
- Reduces coordination overhead across squads

**Example:**
- Notification team adds new feature: WhatsApp alerts in addition to SMS and email
- They extend their service internally to support WhatsApp but keep API (`POST /notifications/send`) the same
- Payment Service continues calling the same API — no changes needed

**Result:** New features delivered faster, without waiting for changes in Payment or Fraud

### 4. Testability

Since APIs are well-defined, services can be tested in isolation using **mocks or stubs**.

**Why it matters:**
- Easier QA and automated testing
- Bugs caught earlier, without needing whole system running

**Example:**
- QA team tests Fraud API by sending fake requests:
  ```json
  { "transactionId": "TXN123", "amount": 5000 }
  ```
- Checks if it returns:
  ```json
  { "fraudulent": false }
  ```
- They don't need Payment Service or Notification Service running at all

**Result:** Testing is faster, cheaper, and less dependent on other services

### 5. Replaceability

A service can be **swapped out or replaced** as long as it continues to respect its API contract.

**Why it matters:**
- Future-proofs the architecture
- Avoids vendor lock-in
- Makes migrations safer

**Example:**
- Notification Service initially built in-house
- Later, FinPay decides to use third-party provider (e.g., Twilio) for SMS/WhatsApp
- As long as new provider accepts same input contract (to, message) and returns same output (status: DELIVERED), Payment Service doesn't change at all

**Result:** Businesses can modernize or change vendors without rewriting entire system

## API Contract Design Principles

### 1. Versioning

Always version APIs to allow evolution:

```
POST /v1/payments  (current)
POST /v2/payments  (new features)
```

Old clients continue using v1, new clients use v2.

### 2. Idempotency

Payment APIs should be idempotent (same request twice = same result):

```json
{
  "requestId": "REQ123",
  "fromAccount": "Ravi123",
  "toAccount": "Landlord456",
  "amount": 5000
}
```

If request is sent twice, don't charge Ravi twice.

### 3. Error Handling

Define clear error responses:

```json
{
  "status": "ERROR",
  "code": "INSUFFICIENT_FUNDS",
  "message": "Account balance is ₹2000, cannot transfer ₹5000"
}
```

### 4. Documentation

Document the contract (OpenAPI/Swagger):

```yaml
/payments:
  post:
    summary: Create a payment
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            properties:
              fromAccount:
                type: string
              toAccount:
                type: string
              amount:
                type: number
```

## Real-World Flow (FinPay Example)

### Scenario: Ravi Pays Rent

**Step 1: Client calls Payment API**
```bash
POST /payments
{
  "fromAccount": "Ravi123",
  "toAccount": "Landlord456",
  "amount": 5000
}
```

**Step 2: Payment Service processes**
- Validates request matches PaymentRequest contract
- Processes transaction
- Returns response matching PaymentResponse contract

**Step 3: Payment Service calls Fraud API**
```bash
POST /fraud/check
{
  "transactionId": "TXN789",
  "amount": 5000
}
```

**Step 4: Fraud Service responds**
```json
{
  "fraudulent": false
}
```

**Step 5: Payment Service calls Notification API**
```bash
POST /notifications/send
{
  "to": "Ravi",
  "message": "Payment of ₹5000 successful"
}
```

**Each service:**
- Respects its API contract
- Doesn't know internals of others
- Can be updated independently

## Key Takeaway

**API contracts are the foundation of loose coupling.**

For FinPay Wallet:
- Payment doesn't know **how** Fraud detects suspicious transactions (rules? ML? AI?)
- Payment doesn't know **how** Notification sends alerts (SMS? Email? WhatsApp? Twilio? in-house?)
- Payment only knows **what** to send and **what** to expect back

This **contract-first design** enables:
- Independent development
- Independent deployment
- Independent scaling
- Independent technology choices
- Independent failure handling

**"Contract stability enables system evolution."** — As long as contracts remain stable, everything behind them can change freely.