# 7. Dependency Inversion and Interface-Driven Design

## Concepts

### What Is Dependency Inversion?

In software design, **dependency inversion** is the principle that high-level modules (business logic) should not depend on low-level modules (implementations). Instead, both should depend on **abstractions** (interfaces).

**Without dependency inversion:**
- Payment Service directly depends on a specific Notification library (say, Twilio SMS)
- If you change providers, you must change Payment code

**With dependency inversion:**
- Payment Service only depends on a Notification **interface**
- Twilio, AWS SNS, or any other provider **implements** that interface
- Payment doesn't care which one is used — it only calls the interface

**This keeps systems flexible, testable, and future-proof.**

### What Problem Does It Solve?

**1. Tight Coupling**
- Without interfaces, one service is locked to a particular implementation
- **Example:** Payment tied to Twilio for SMS → hard to switch

**2. Difficult Testing**
- If Payment always uses Twilio, how do you test it without sending real SMS?
- You need a fake/mock provider, which only works with interfaces

**3. Reduced Flexibility**
- Businesses evolve: maybe today you use Twilio, tomorrow you need AWS SNS for cost savings
- Without dependency inversion, this requires rewriting core Payment logic

## FinPay Wallet Example

Let's say FinPay wants to send payment alerts.

### Without Dependency Inversion

```go
func SendPaymentNotification(msg string) {
    // Hardcoded dependency
    twilio.SendSMS(msg)
}
```

**Problem:** Payment is married to Twilio. If you change to AWS SNS, you rewrite this function.

### With Dependency Inversion

```go
type Notifier interface {
    Send(message string) error
}
```

**Benefits:**
- Payment depends only on `Notifier`
- Twilio implements `Notifier`
- AWS SNS can also implement `Notifier`
- Payment doesn't care who's behind the scenes

**Result:** Payment logic is stable, while notification providers can change freely

## Prerequisites

- Go 1.20+
- Understanding of interfaces
- Basic understanding of dependency injection

## How to Run

### Step 1: Run the Interface Demo

```bash
cd chapter8-loose-coupling/07-dependency-inversion

# Run the demo
go run interface_demo.go
```

**Expected output:**
```
Payment processed successfully
Twilio SMS sent: Payment successful notification
Payment processed successfully
AWS SNS sent: Payment successful notification
```

**Observation:**
- Same `processPayment` function works with different implementations
- Payment logic never changes
- Can switch providers at runtime

### Step 2: Understand the Code

**1. Define the Abstraction (Interface):**
```go
type Notifier interface {
    Send(message string) error
}
```

This is the **contract**. Payment only knows this.

**2. Implement Twilio:**
```go
type TwilioNotifier struct{}

func (t TwilioNotifier) Send(message string) error {
    fmt.Println("Twilio SMS sent:", message)
    return nil
}
```

**3. Implement AWS SNS:**
```go
type SNSNotifier struct{}

func (s SNSNotifier) Send(message string) error {
    fmt.Println("AWS SNS sent:", message)
    return nil
}
```

**4. Payment Depends Only on Interface:**
```go
func processPayment(n Notifier) {
    fmt.Println("Payment processed successfully")
    n.Send("Payment successful notification")
}
```

**Key:** `processPayment` doesn't know (or care) if it's Twilio or AWS — it just calls `Notifier.Send`.

**5. Runtime Choice:**
```go
func main() {
    twilio := TwilioNotifier{}
    aws := SNSNotifier{}

    processPayment(twilio)  // uses Twilio
    processPayment(aws)     // uses AWS SNS
}
```

At runtime, we choose which implementation to inject.

**This demonstrates dependency inversion in action:** high-level Payment logic is independent of low-level Notification details.

## Dependency Inversion Principle (DIP)

### Traditional Dependency Flow

```
┌──────────────────┐
│  Payment Logic   │
│  (High-level)    │
└────────┬─────────┘
         │ depends on
         ▼
┌──────────────────┐
│  Twilio SDK      │
│  (Low-level)     │
└──────────────────┘
```

**Problem:** Change Twilio → Must change Payment

### Inverted Dependency Flow

```
┌──────────────────┐
│  Payment Logic   │
│  (High-level)    │
└────────┬─────────┘
         │ depends on
         ▼
┌──────────────────┐
│  Notifier        │
│  (Interface)     │
└────────┬─────────┘
         ▲ implements
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼────┐
│Twilio │ │AWS SNS│
└───────┘ └───────┘
```

**Solution:** Both Payment and Twilio depend on Notifier interface

## Benefits

### 1. Flexibility
- Swap providers (Twilio → AWS SNS) without rewriting Payment
- **Example:** FinPay starts with Twilio, later switches to AWS SNS for cost savings
- Only change: which implementation is injected
- Payment code: **zero changes**

### 2. Testability
- Use a fake `MockNotifier` in tests instead of real SMS providers
- **Example:**

```go
type MockNotifier struct {
    MessagesSent []string
}

func (m *MockNotifier) Send(message string) error {
    m.MessagesSent = append(m.MessagesSent, message)
    return nil
}

func TestPaymentFlow(t *testing.T) {
    mock := &MockNotifier{}
    processPayment(mock)

    if len(mock.MessagesSent) != 1 {
        t.Error("Expected 1 notification")
    }
}
```

**No real SMS sent, no API costs, fast tests**

### 3. Future-Proofing
- New providers can be added easily, just by implementing the interface
- **Example:** FinPay adds WhatsApp Business

```go
type WhatsAppNotifier struct{}

func (w WhatsAppNotifier) Send(message string) error {
    fmt.Println("WhatsApp sent:", message)
    return nil
}
```

**Payment code:** still unchanged

### 4. Cleaner Code
- Payment focuses only on business rules, not technical details
- **Example:** Payment doesn't know about Twilio API keys, AWS regions, or HTTP endpoints
- All that complexity is hidden behind the interface

## Real-World Application

### Configuration-Based Provider Selection

```go
func main() {
    provider := os.Getenv("NOTIFICATION_PROVIDER")

    var notifier Notifier

    switch provider {
    case "twilio":
        notifier = TwilioNotifier{}
    case "aws":
        notifier = SNSNotifier{}
    case "whatsapp":
        notifier = WhatsAppNotifier{}
    default:
        notifier = TwilioNotifier{} // default
    }

    processPayment(notifier)
}
```

**Same code, different environments:**
- Dev: Uses MockNotifier (no real SMS)
- Staging: Uses Twilio
- Production: Uses AWS SNS

**Change provider without code deployment — just update environment variable**

## SOLID Principles Connection

Dependency Inversion is the "D" in SOLID:
- **S**ingle Responsibility
- **O**pen/Closed
- **L**iskov Substitution
- **I**nterface Segregation
- **D**ependency Inversion

### How It Relates to Loose Coupling

**Dependency Inversion enables loose coupling:**
- Services depend on interfaces (contracts)
- Not on concrete implementations (code)
- Can swap implementations without touching business logic

## Interface Design Best Practices

### 1. Keep Interfaces Small

```go
// ✅ Good: Small, focused interface
type Notifier interface {
    Send(message string) error
}

// ❌ Bad: Large, complex interface
type NotificationSystem interface {
    Send(message string) error
    Schedule(message string, time time.Time) error
    Cancel(messageId string) error
    GetStatus(messageId string) (string, error)
    UpdatePreferences(userId string, prefs map[string]bool) error
}
```

**Principle:** "Interface segregation" — clients shouldn't depend on methods they don't use

### 2. Name Interfaces After Behavior

```go
// ✅ Good: Describes what it does
type Sender interface { ... }
type Validator interface { ... }
type Processor interface { ... }

// ❌ Bad: Generic names
type Service interface { ... }
type Manager interface { ... }
type Handler interface { ... }
```

### 3. Accept Interfaces, Return Structs

```go
// ✅ Good: Accept interface
func processPayment(n Notifier) { ... }

// ❌ Bad: Accept concrete type
func processPayment(n TwilioNotifier) { ... }
```

**Principle:** Be liberal in what you accept, conservative in what you return

## Key Takeaway

**Dependency inversion is the foundation of flexible, testable, maintainable systems.**

For FinPay Wallet:
- **Without dependency inversion:** Payment locked to Twilio → hard to test, hard to change
- **With dependency inversion:** Payment depends on Notifier interface → easy to test, easy to swap providers

**Design principle:** "Depend on abstractions, not concretions."

Benefits:
- ✅ Flexible — swap implementations without touching business logic
- ✅ Testable — use mocks for fast, reliable tests
- ✅ Future-proof — add new providers easily
- ✅ Clean — business logic stays focused on business rules

**"High-level policy should not depend on low-level details. Both should depend on abstractions."** — Robert C. Martin (Uncle Bob)