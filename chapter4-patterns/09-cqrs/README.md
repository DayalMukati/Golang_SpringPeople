# CQRS (Command Query Responsibility Segregation) Example

This example demonstrates the CQRS pattern by separating write operations (commands) from read operations (queries).

## Concept

CQRS separates how you write data (commands) from how you read data (queries):

- **Commands** = requests that change state (create payment, update KYC, transfer funds)
  - Validated, processed, may trigger events

- **Queries** = requests that fetch data (show balance, list transactions)
  - Never change state, just return info

## Implementation

### Write Model (Ledger)
- Stores all payment transactions
- Each payment has: UserID, Amount, Type (debit/credit), TxID

### Read Model (Balances)
- Pre-computed balance view for fast queries
- Updated asynchronously when commands are processed

## Endpoints

### Command: POST /command/pay
Process a payment (write side)

```bash
curl -X POST http://localhost:9000/command/pay \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":100,"type":"credit","tx_id":"tx1"}'
```

### Query: GET /query/balance?user=U1
Get user balance (read side)

```bash
curl "http://localhost:9000/query/balance?user=U1"
```

## How to Run

1. Start the server:
```bash
go run main.go
```

2. Credit $100 to U1:
```bash
curl -X POST http://localhost:9000/command/pay \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":100,"type":"credit","tx_id":"tx1"}'
```

3. Debit $30 from U1:
```bash
curl -X POST http://localhost:9000/command/pay \
  -H "Content-Type: application/json" \
  -d '{"user_id":"U1","amount":30,"type":"debit","tx_id":"tx2"}'
```

4. Query balance:
```bash
curl "http://localhost:9000/query/balance?user=U1"
# Output: {"user":"U1","balance":70.00}
```

## Benefits

- **Performance**: Reads optimized separately from writes
- **Scalability**: Read and write sides can scale independently
- **Flexibility**: Different models for different needs
- **Clarity**: Clear separation of concerns