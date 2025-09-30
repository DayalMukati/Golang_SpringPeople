# 6. Caching for Performance and Scalability

## Concepts

### What is Caching?

**Caching** means storing frequently accessed data in a fast, temporary storage layer (memory or in-memory databases like Redis).

Instead of fetching Ravi's balance from the main database every time, the system can serve it from a cache if it hasn't changed. This reduces load on the DB and makes responses nearly instant.

**Analogy:** Cache is like keeping a copy on your desk for quick access, instead of walking to the library every time you need a fact.

### Why Caching Matters

Databases are powerful but slow compared to memory. If every balance check, fraud detection, and transaction history query hits the DB, it will eventually choke.

**Caching solves:**
1. **High Read Traffic** → Moves repetitive queries out of the database
2. **Latency Issues** → Returns results from memory in milliseconds
3. **Scalability Bottlenecks** → Lets databases focus on critical writes while cache handles reads

### FinPay Wallet Examples

**1. Balance Checks**
- Ravi opens his FinPay app 10 times a day just to check balance
- Without caching → every request hits the DB
- With caching → balance stored in Redis for 1 minute; next 9 checks come from cache

**2. Fraud Detection**
- Fraud rules often need reference data (blacklisted accounts, geo-locations)
- Keeping this data in cache avoids hitting DB for every transaction

**3. Transaction History**
- When Ravi views last 10 transactions, system caches result for 30 seconds
- Even if 100K users open "Transaction History" at once, most requests come from cache, not DB

**Result:** System feels fast for Ravi while DB stays healthy under peak load

## Prerequisites

- Go 1.20+
- Basic understanding of maps
- Redis (optional, for production examples)

## How to Run

### Step 1: Run Simple Cache Demo

```bash
cd chapter7-scalability/06-caching

# Run the demo
go run cache_demo.go
```

**Expected output:**
```
Cache miss for Ravi - fetching from DB...
Balance: 25000
Cache hit for Ravi
Balance: 25000
Cache hit for Ravi
Balance: 25000
```

### Step 2: Understand the Code

**Explanation:**
- **First call** → "Cache miss" → hits DB, stores value in cache
- **Next calls** → "Cache hit" → retrieves from cache (much faster)

In production, Redis or Memcached would replace this simple map.

### Step 3: Production Caching with Redis (Conceptual)

**Install Redis:**
```bash
# macOS
brew install redis
redis-server

# Linux
sudo apt-get install redis-server
sudo systemctl start redis

# Docker
docker run -d -p 6379:6379 redis:latest
```

**Go code with Redis:**
```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

func getBalance(userId string) (int, error) {
    // Try cache first
    val, err := rdb.Get(ctx, "balance:"+userId).Int()
    if err == nil {
        fmt.Println("Cache hit")
        return val, nil
    }

    // Cache miss - fetch from DB
    fmt.Println("Cache miss - querying DB")
    balance := 25000 // Simulate DB query

    // Store in cache with 60-second expiry
    rdb.Set(ctx, "balance:"+userId, balance, 60*time.Second)

    return balance, nil
}

func main() {
    balance, _ := getBalance("ravi")
    fmt.Println("Balance:", balance)

    // Second call hits cache
    balance, _ = getBalance("ravi")
    fmt.Println("Balance:", balance)
}
```

**Dependencies:**
```bash
go get github.com/go-redis/redis/v8
```

## Data Retrieval Process

### Without Caching
```
User Request → API Server → Database → API Server → User
                    ↓                      ↑
                  50-100ms             Total: 50-100ms
```

**Every request:**
- Hits database (50-100ms latency)
- Database handles 100K requests/minute
- Database becomes bottleneck

### With Caching
```
User Request → API Server → Cache (Redis) → API Server → User
                    ↓              ↑             ↑
                  1-2ms      Cache Hit    Total: 1-2ms

           [On Cache Miss]
              ↓
         Database (50-100ms)
              ↓
         Update Cache
              ↓
         Return to User
```

**Most requests:**
- Hit cache (1-2ms latency) — **50× faster**
- Only cache misses hit database
- Database handles maybe 5K requests/minute instead of 100K

## Benefits of Caching in Fintech

### 1. Performance Boost
- Fetching data from memory is **orders of magnitude faster** than querying database on disk
- **Example:** FinPay Wallet retrieves Ravi's balance from Redis in <1ms, compared to 50-100ms from database
- During peak hours (midnight salary credits), millions of users get instant responses
- **Impact:** Customers feel app is snappy and reliable, building trust

### 2. Reduced Database Load
- Every redundant read avoided is one less burden on database
- **Example:** 1 million users open app within 10 minutes
  - Without caching: 1M DB queries
  - With caching: 50K DB queries (only first-time lookups, cache misses)
- Database free to focus on critical writes — debits, credits, fraud checks
- **Prevents:** Slowdowns and outages caused by overloaded databases

### 3. Scalability
- By offloading reads to cache, system effectively increases throughput capacity without changing database
- **Example:** FinPay handles 200,000 balance checks per minute on salary day
- With caching, DB only sees fraction of those
- **Makes:** System scale to millions of concurrent users without costly DB overhauls

### 4. Cost Savings
- Scaling databases is expensive (larger instances, SSDs, HA replicas)
- Scaling caches is much cheaper
- **Example:** Instead of upgrading primary database to 64-core enterprise server, FinPay adds a few Redis nodes for fraction of cost
- Cloud providers charge less for memory-based cache clusters than for top-tier DB instances
- **Helps:** Fintech companies deliver enterprise-grade performance on startup budgets

### 5. Improved Reliability in Spikes
- During sudden surges, cached responses absorb most traffic
- **Example:** On salary day, 80% of balance checks come from cache
- Prevents database from tipping over under extreme load
- Even if DB gets temporarily slow, cached results buy breathing room
- **Users:** Continue to experience fast responses even when backend systems strain

### 6. Better User Experience for Derived Data
- Some data doesn't change often but is expensive to compute (fraud risk scores, monthly statements)
- **Example:** Ravi opens "Last 10 Transactions" screen
- System caches it for 30 seconds
- If he reopens within that window, cache serves result instantly
- **Users:** Experience consistent performance, even on features that would otherwise be slow

## Caching Strategies

### 1. Cache-Aside (Lazy Loading)
**How it works:**
1. Check cache
2. If miss, query database
3. Store result in cache
4. Return to user

**Best for:** Read-heavy, infrequently changing data (user profiles, balances)

**Example:** FinPay balance checks

### 2. Write-Through
**How it works:**
1. Write to database
2. Immediately update cache
3. Return success

**Best for:** Data that must be consistent (account balances after transaction)

**Example:** After Ravi sends ₹5,000, update both DB and cache

### 3. Write-Behind (Write-Back)
**How it works:**
1. Write to cache immediately
2. Asynchronously write to database later
3. Return success quickly

**Best for:** High write throughput (logging, analytics)

**Risk:** Data loss if cache crashes before DB write

### 4. Refresh-Ahead
**How it works:**
1. Predict which data will be needed
2. Pre-load into cache before expiry
3. Users always hit warm cache

**Best for:** Predictable access patterns (morning login surge)

## Cache Expiration Strategies

### Time-Based (TTL)
- Set expiry time when storing in cache
- **Example:** Balance cached for 60 seconds

```go
rdb.Set(ctx, "balance:ravi", 25000, 60*time.Second)
```

### Event-Based
- Invalidate cache when data changes
- **Example:** When transaction completes, delete old balance from cache

```go
// After successful transaction
rdb.Del(ctx, "balance:ravi")
```

### LRU (Least Recently Used)
- Cache evicts least recently used items when full
- **Example:** Redis with maxmemory-policy=allkeys-lru

## What to Cache in Fintech?

### ✅ Good Candidates
- **User profiles** (name, email, KYC status) — rarely change
- **Current balances** — frequently read, cache for 10-60 seconds
- **Transaction history** — recent 10-50 transactions, cache for 30 seconds
- **Fraud rules** — blacklisted accounts, geo-restrictions
- **Exchange rates** — update every 5 minutes
- **Product catalog** — credit cards, loan offers

### ❌ Bad Candidates
- **Active transactions in-flight** — too critical, risk of inconsistency
- **Real-time fraud scores** — must be computed fresh each time
- **Regulatory audit logs** — must be 100% accurate, can't risk cache corruption
- **PII (Personally Identifiable Information)** — security risk if cache is compromised

## Real-World Flow (FinPay Example)

### Scenario: Salary Day Balance Checks

**00:00 - Surge begins**
- 1 million users open app to check if salary credited
- Average user checks balance 3 times in first minute

**Without caching:**
- Database receives: 3M queries in 1 minute (50K/second)
- DB CPU: 100% (maxed out)
- Query latency: 5 seconds (users frustrated)
- Some queries timeout
- **Result:** Poor user experience, potential outage

**With caching (60-second TTL):**
- First check per user: Cache miss → DB query → Store in cache
- Second check: Cache hit (instant)
- Third check: Cache hit (instant)
- Database receives: 1M queries in 1 minute (16K/second) — **67% reduction**
- DB CPU: 40% (healthy)
- Cache latency: <1ms
- **Result:** Smooth user experience, system stable

**Cost impact:**
- Without cache: Need to scale DB to handle 50K/second → $5K/month
- With cache: DB handles 16K/second + Redis cluster → $1.5K/month
- **Savings:** $3.5K/month

## Monitoring Cache Performance

### Key Metrics

**1. Cache Hit Ratio**
```
Hit Ratio = Cache Hits / (Cache Hits + Cache Misses)
```
- **Target:** >90% for read-heavy workloads
- **Example:** 95K cache hits, 5K cache misses → 95% hit ratio ✅

**2. Cache Latency**
- **Target:** <5ms for in-memory cache
- **Example:** Redis GET takes 1.2ms ✅

**3. Eviction Rate**
- How often cache is full and evicting items
- High eviction rate → need more cache capacity

**4. Memory Usage**
- Keep below 80% to avoid performance degradation

### Redis Monitoring Commands
```bash
# Get cache stats
redis-cli INFO stats

# Check memory usage
redis-cli INFO memory

# Monitor real-time commands
redis-cli MONITOR
```

## Challenges and Solutions

### Challenge 1: Cache Invalidation
**Problem:** "There are only two hard things in Computer Science: cache invalidation and naming things" — Phil Karlton

**Solution:**
- Use short TTL for frequently changing data (10-60 seconds)
- Event-based invalidation for critical updates
- Version keys to prevent stale data

### Challenge 2: Cache Stampede
**Problem:** Cache expires, 10K requests simultaneously hit DB

**Solution:**
- Use locking: first request fetches from DB, others wait
- Stagger TTLs with random jitter
- Implement refresh-ahead pattern

### Challenge 3: Data Consistency
**Problem:** Cache has ₹25,000, DB has ₹20,000 (Ravi just spent ₹5,000)

**Solution:**
- Invalidate cache on write
- Use write-through for critical data
- Accept eventual consistency for non-critical reads

### Challenge 4: Cache Penetration
**Problem:** Queries for non-existent data bypass cache, hit DB every time

**Solution:**
- Cache "null" results with short TTL
- Bloom filters to check existence before querying

## Best Practices

1. **Cache at multiple layers**: Browser → CDN → API → Database
2. **Set appropriate TTLs**: Balance freshness vs hit ratio
3. **Monitor hit ratio**: Should be >80%, ideally >90%
4. **Use Redis Cluster**: For high availability and scalability
5. **Compress large values**: Save memory and network bandwidth
6. **Implement circuit breakers**: Fallback to DB if cache fails
7. **Warm up cache**: Pre-load popular data after deployment
8. **Security**: Encrypt sensitive data in cache, use authentication

## Key Takeaway

Caching is one of the **highest-impact, lowest-effort optimizations** for cloud-native fintech systems.

For FinPay Wallet:
- **50-100× faster response times** for cached data
- **80-95% reduction in database load** during peak traffic
- **Significant cost savings** compared to scaling databases
- **Better user experience** during salary day surges

Combined with auto-scaling and load balancing, caching completes the scalability triad:
- **Load balancing** distributes requests across pods
- **Auto-scaling** adds capacity when needed
- **Caching** reduces backend load and improves speed

Together, these enable FinPay Wallet to deliver bank-grade performance and reliability while controlling costs — essential for fintech startups competing with established players.