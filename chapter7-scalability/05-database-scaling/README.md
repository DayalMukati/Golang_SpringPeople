# 5. Database Scaling Approaches (Sharding, Replication)

## Concepts

Scaling application servers with pods is straightforward, but **databases are harder**. They store state (user balances, transactions, audit logs), and you can't simply spin up 30 new copies without thinking about consistency and data integrity.

In fintech, the database must:
- Handle spikes in writes and reads (salary day transactions)
- Ensure accuracy (no double-debits, no missing transactions)
- Provide availability (can't go offline when a node crashes)

**Two core scaling techniques:** Replication and Sharding

## Approach 1: Replication

### What Is Replication?

**Replication** means making copies of the same database across multiple nodes.

**Architecture:**
- **Primary (Master):** Handles all writes (INSERT, UPDATE, DELETE)
- **Replicas (Slaves):** Handle reads (SELECT queries)

This spreads the load: write-heavy operations go to one node, while read-heavy queries spread across many replicas.

### FinPay Wallet Example

- **Primary DB:** Processes wallet updates when Ravi sends money to his landlord
- **Read Replicas:** Handle balance lookups, fraud checks, reporting dashboards
- When Ravi's app shows "Your balance: ₹25,000," that's likely served from a read replica

### Pros
✅ Easy to set up (supported by MySQL, Postgres, MongoDB, etc.)
✅ Improves read throughput dramatically
✅ Replicas provide failover if primary goes down

### Cons
❌ Replication lag → new balance may take milliseconds/seconds to appear on replicas
❌ Writes are still limited to single primary

### Benefits of Replication in Fintech

**1. Faster Reads**
- Most fintech apps are read-heavy
- Every time Ravi opens FinPay Wallet, app checks balance and fetches last 10 transactions
- With replication, these requests hit read replicas, leaving primary free to handle money movements
- **Result:** Balance checks and dashboards stay fast, even during peak loads

**2. Reliability and Failover**
- If primary database crashes at midnight on salary day, one replica can be promoted to take over writes
- Users still access balances and histories while system repairs itself
- **Ensures:** Business continuity, avoids downtime

**3. Scalability for Analytics**
- Fintech apps need fraud detection, compliance reports, insights dashboards
- Running heavy queries on primary would slow down transactions
- Replicas offload analytics queries without affecting live payments
- **Keeps:** Payments fast while meeting regulatory reporting needs

## Approach 2: Sharding

### What Is Sharding?

**Sharding** means splitting the database by data range or key and distributing across multiple nodes. Each shard holds a subset of the data.

### FinPay Wallet Example

**Sharding strategy:**
- Users A–M stored on Shard 1
- Users N–Z stored on Shard 2
- Each shard handles both reads and writes, but only for its assigned users

When Ravi queries his balance, the app knows he belongs to Shard 2 and routes the request there.

### Pros
✅ Solves both read and write scaling
✅ Near-unlimited growth — add more shards as user base grows

### Cons
❌ Complexity in managing queries (app must know which shard to hit)
❌ Cross-shard queries (e.g., "top 10 users across all shards") are harder

### Benefits of Sharding in Fintech

**1. Handles Write-Heavy Workloads**
- On salary day, millions of users transfer money simultaneously
- If all writes go to one DB, it becomes a bottleneck
- With sharding:
  - Ravi's transaction (user ID in shard 2) → Shard 2
  - Meena's transaction (user ID in shard 1) → Shard 1
- Each shard handles smaller piece of workload
- **System can process more payments in parallel**

**2. Near-Infinite Growth**
- Single DB server has limits (CPU, memory, disk)
- With sharding, FinPay can keep adding shards as customers grow
- When FinPay expands from 1M to 10M users, add Shard 3 and Shard 4
- **Growth becomes predictable and unbounded**

**3. Better Locality of Data**
- In global fintech apps, shards can be placed near customers
- **Example:** EU users on shard in Frankfurt, Asia users on shard in Singapore
- **Benefits:** Reduces latency, helps with data residency compliance

## Combined Approach: Replication + Sharding

In practice, fintech systems like FinPay use **both approaches**:
- **Sharding** splits data so no single DB is overloaded
- **Replication** inside each shard ensures high availability and faster reads

### Example Architecture

**Ravi's account is on Shard 2:**
- Shard 2 has 1 primary DB (for writes)
- Shard 2 has 3 replicas (for reads, analytics, failover)
- Ravi's balance check → hits replica
- Ravi's rent payment → hits primary
- Compliance reports → query replicas, leaving primary free

**This combined model ensures:**
- Performance (load distributed)
- Reliability (multiple copies, failover)
- Regulatory compliance (audit queries don't slow payments)

## Prerequisites

- Go 1.20+
- Basic understanding of maps and data structures
- Conceptual understanding of databases (MySQL, Postgres, MongoDB)

## How to Run

### Step 1: Understand Sharding Concept

```bash
cd chapter7-scalability/05-database-scaling

# Run the sharding demo
go run sharding_demo.go
```

**Expected output:**
```
Single DB (no sharding):
map[Alice:1000 Meena:500 Ravi:25000]

Sharded DBs:
Shard 1: map[Alice:1000]
Shard 2: map[Meena:500 Ravi:25000]

User 'Ravi' belongs to Shard 2
Balance: 25000 (from Shard 2)
```

### Step 2: Understanding the Code

The demo shows how data is split:

**Single map (one DB):**
- All users stored together
- As it grows too big, becomes slow

**Sharded maps:**
- Split into multiple maps (shard1, shard2)
- Queries go only to shard that holds the user
- Reduces load per database server

### Step 3: Replication in Real Databases

**MySQL Replication Example (Conceptual):**

```bash
# On Primary (writes)
mysql> INSERT INTO wallets (user_id, balance) VALUES ('ravi', 25000);

# Replicas automatically receive the change
# On Replica 1, 2, 3 (reads)
mysql> SELECT balance FROM wallets WHERE user_id = 'ravi';
```

**Application code routes queries:**
```go
// Write operations → Primary
db.Primary.Exec("UPDATE wallets SET balance = ? WHERE user_id = ?", newBalance, userId)

// Read operations → Replica
db.Replica.Query("SELECT balance FROM wallets WHERE user_id = ?", userId)
```

### Step 4: Sharding in Real Databases

**Application must determine shard:**

```go
func getShardForUser(userId string) int {
    // Hash-based sharding
    hash := hashCode(userId)
    return hash % totalShards
}

func getBalance(userId string) (int, error) {
    shard := getShardForUser(userId)
    db := shardConnections[shard]

    var balance int
    err := db.QueryRow("SELECT balance FROM wallets WHERE user_id = ?", userId).Scan(&balance)
    return balance, err
}
```

**Database tools that help:**
- **MySQL:** Use Vitess for automatic sharding
- **MongoDB:** Built-in sharding support
- **PostgreSQL:** Use Citus extension for sharding
- **Cassandra:** Automatic sharding by partition key

## Replication vs Sharding: When to Use What?

### Use Replication When:
- **Read-heavy workload** (e.g., balance checks, transaction history)
- **Relatively small dataset** (fits on one machine)
- **Need high availability** (failover capability)

**Example:** FinPay Dashboard service queries user data for analytics → use read replicas

### Use Sharding When:
- **Write-heavy workload** (e.g., transaction processing)
- **Large dataset** (doesn't fit on one machine)
- **Need horizontal write scaling**

**Example:** FinPay Transaction service processing millions of payments → use sharding

### Use Both When:
- **Very large scale** (millions of users, high read + write load)
- **Global deployment** (data locality requirements)
- **Critical availability** (no single point of failure)

**Example:** FinPay production system serving 10M users globally

## Real-World Flow (FinPay Example)

### Without Scaling

**Salary Day:**
- Single database handles 200,000 transactions/hour
- DB CPU: 100% (maxed out)
- Write latency: 5 seconds per transaction
- Read latency: 3 seconds per query
- **Result:** System crashes, payments fail

### With Replication Only

**Salary Day:**
- Primary handles 200,000 writes/hour
- 5 replicas handle reads
- Write latency: Still 4 seconds (primary overloaded)
- Read latency: 50ms (replicas healthy)
- **Result:** Reads are fast, but writes still bottleneck

### With Sharding Only

**Salary Day:**
- 4 shards each handle 50,000 transactions/hour
- Write latency: 200ms (load distributed)
- Read latency: 500ms (no read optimization)
- **Result:** Writes are fast, but reads could be better

### With Replication + Sharding (Best)

**Salary Day:**
- 4 shards (each has 1 primary + 3 replicas)
- Each shard primary handles 50,000 writes/hour
- Each shard's 3 replicas handle reads
- Write latency: 150ms (distributed writes)
- Read latency: 50ms (distributed reads)
- **Result:** Both reads and writes are fast and scalable

## Monitoring Database Scaling

### Key Metrics to Watch

**Replication:**
- Replication lag (should be < 1 second)
- Replica CPU/memory usage
- Failed replica count

**Sharding:**
- Shard distribution (even vs skewed)
- Cross-shard query count (minimize these)
- Shard rebalancing events

**Commands (MySQL example):**
```sql
-- Check replication lag
SHOW SLAVE STATUS\G

-- Check shard distribution
SELECT shard_id, COUNT(*) FROM users GROUP BY shard_id;
```

## Challenges and Solutions

### Challenge 1: Replication Lag
**Problem:** Ravi transfers money, but balance doesn't update immediately on replica

**Solution:**
- For critical reads (balance after transfer), read from primary
- For non-critical reads (dashboard), read from replica
- Use eventual consistency model

### Challenge 2: Cross-Shard Queries
**Problem:** "Show me top 10 spenders across all shards" requires querying all shards

**Solution:**
- Aggregate at application layer
- Use separate analytics database (data warehouse)
- Pre-compute common aggregations

### Challenge 3: Shard Rebalancing
**Problem:** Shard 1 has 10M users, Shard 2 has 1M users (uneven)

**Solution:**
- Use consistent hashing for better distribution
- Plan for rebalancing during low-traffic windows
- Use tools like Vitess that handle rebalancing automatically

### Challenge 4: Transaction Spanning Shards
**Problem:** Transfer money from Alice (Shard 1) to Ravi (Shard 2) requires two-phase commit

**Solution:**
- Design to avoid cross-shard transactions when possible
- Use saga pattern for distributed transactions
- Implement idempotency for retry safety

## Best Practices

1. **Start with replication** for read scaling
2. **Add sharding** when write load becomes bottleneck
3. **Shard by user ID** in fintech (keeps user data together)
4. **Monitor replication lag** closely
5. **Plan shard key carefully** (hard to change later)
6. **Use managed services** (AWS RDS read replicas, MongoDB Atlas sharding) when possible
7. **Test failover** regularly (chaos engineering)

## Key Takeaway

Database scaling is **essential but complex** for fintech systems:

- **Replication** solves read scalability and provides high availability
- **Sharding** solves write scalability and enables unlimited growth
- **Together**, they provide the foundation for systems that can:
  - Handle millions of concurrent users
  - Process hundreds of thousands of transactions per minute
  - Maintain low latency even during peak load
  - Survive database failures without downtime

For FinPay Wallet, this means Ravi's ₹5,000 rent payment goes through smoothly on salary day midnight, just as it would on a quiet Tuesday afternoon — regardless of the millions of other transactions happening simultaneously.