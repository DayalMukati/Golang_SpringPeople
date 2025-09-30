# 9. Cost Considerations in Scaling Architectures

## Concepts

### What Is This About?

**"Cost considerations"** means designing your scaling approach so the system stays fast and the cloud bill stays sane.

In practice, you balance three forces:
- **Performance** — meet SLOs at peak
- **Reliability** — survive failures, multi-AZ/region
- **Cost** — pay only for value, not waste

In cloud-native systems you don't buy servers once—you **rent capacity continuously**. Every design choice (pods, replicas, shards, regions, logs) turns into a line item on the invoice.

### What Problem Does It Solve?

**Without a cost-aware design, teams either:**
- **Over-provision for worst case** → great performance, runaway spend
- **Under-provision to save money** → timeouts, failed payments, reputational damage

**A cost-informed architecture meets demand precisely and scales down aggressively when idle.**

## Fintech Example: FinPay Wallet

### Context

**Normal load:** ~10k transactions/hour
**Salary day midnight spike:** ~200k transactions/hour for ~3 hours

**Core services:**
- Payment API (stateless)
- Fraud (CPU-heavy)
- Wallet (DB writes)
- Notification (I/O heavy)
- Redis cache
- Postgres primary + 2 replicas

### Two Strategies

**Strategy 1: Always-On Peak Capacity**
- 30 Payment pods 24×7
- 8 Fraud pods 24×7
- RDS at top tier
- **Pros:** Simple
- **Cons:** Pay peak rates all month

**Strategy 2: Elastic Baseline + Burst (Recommended)**
- Payment: 5→30 pods via HPA (min=5, max=30)
- Fraud: 2→8 pods
- Postgres: right-sized primary + replicas
- Hot Redis cache reduces DB reads
- Same-AZ colocation to avoid cross-AZ traffic for chatty services
- Observability retention tuned (90 days metrics, 7 days full logs)

**Result:** Salary day is smooth, off-peak cost drops sharply

## The Big Cost Drivers (and How to Steer Them)

### 1. Compute (Pods / Nodes)

**Pricing tiers:**
- **On-demand** — most expensive, instant availability
- **Reserved/Savings Plans** — ~40-60% discount, 1-3 year commitment
- **Spot/Preemptible** — ~70-90% discount, can be evicted

**Tactics:**
- Baseline on reserved nodes
- Burst on on-demand or spot for stateless services
- Set HPA min/max and cool-downs to prevent thrash

**Example:**
```yaml
# Baseline on reserved
minReplicas: 5  # Runs 24×7 on reserved instances

# Burst on on-demand
maxReplicas: 30  # Extra 25 pods on on-demand during peaks
```

### 2. Storage (DB, Block, Object)

**Cost factors:**
- **Provisioned IOPS DBs** — very expensive
- **General Purpose SSD** — balanced cost/performance
- **Cold storage** — cheap for archives

**Tactics:**
- Use Redis to offload reads from expensive DB
- Archive old reports to object storage (S3/GCS)
- Choose storage classes by access pattern (hot vs cold)

**Example:**
```
Hot data (last 30 days): Postgres + Redis cache
Warm data (30-90 days): Postgres compressed tables
Cold data (>90 days): S3 Glacier for compliance
```

### 3. Network (Egress)

**Cost factors:**
- **Cross-AZ traffic** — $0.01-0.02 per GB
- **Internet egress** — $0.08-0.12 per GB
- **Same-AZ traffic** — FREE

**Tactics:**
- Keep chatty services in one AZ
- Place shards near users (geo-distribution)
- Compress payloads
- Batch writes

**Example:**
```
Payment ←→ Fraud ←→ Wallet
    ↓
All in same AZ (us-east-1a)
= Zero network cost

Payment (us-east-1a) ←→ Fraud (us-west-2)
= $0.02/GB × 1TB/month = $20/month extra
```

### 4. Managed Services

**Cost factors:**
- **Kafka/Redis/RDS** tiers grow with throughput/retention
- **Retention policies** directly impact cost

**Tactics:**
- Enforce quotas and TTLs
- Topic/stream retention limits
- Only retain high-value data at high fidelity

**Example:**
```
Kafka topics:
- transactions: 7-day retention (compliance)
- logs: 1-day retention (debugging)
- analytics: 90-day retention (business intelligence)
```

### 5. Observability

**Cost factors:**
- **Debug-level logs** — 10-100× more volume than INFO
- **Long retention** — linear cost growth
- **Full trace sampling** — expensive at scale

**Tactics:**
- Sample traces (1-10% of requests)
- Downsample metrics (1min → 5min → 1hour aggregations)
- Set per-service log budgets
- Use structured logging for efficient queries

**Example:**
```
Production logs:
- Level: INFO (not DEBUG)
- Retention: 7 days
- Sampling: 1% for traces

Cost: $200/month

Versus:
- Level: DEBUG
- Retention: 90 days
- Sampling: 100%

Cost: $8,000/month ❌
```

## Concrete Techniques That Save Money (Without Hurting SLOs)

### 1. Right-Size Requests/Limits
- Stop giving every pod 1 CPU "just in case"
- Review with usage data weekly
- Consider VPA recommendations

**Example:**
```yaml
# Before (wasteful)
resources:
  requests:
    cpu: 1000m
    memory: 2Gi

# After (right-sized)
resources:
  requests:
    cpu: 250m      # Actual usage: 180m p95
    memory: 512Mi  # Actual usage: 380Mi p95
```

**Savings:** 75% compute cost reduction per pod

### 2. Autoscaling Guardrails
- HPA min/max
- Scale-up/down policies
- Cooldowns

**Example:**
```yaml
behavior:
  scaleDown:
    stabilizationWindowSeconds: 300  # Wait 5min before scaling down
    policies:
    - type: Percent
      value: 50     # Scale down max 50% at a time
      periodSeconds: 60
```

**Prevents:** Thrashing that wastes money and destabilizes system

### 3. Scale-to-Zero for Spiky Workloads
- Use jobs, event triggers, or KEDA
- Batch/reporting, webhooks, non-critical components

**Example:**
```yaml
# KEDA ScaledObject
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: report-generator
spec:
  minReplicaCount: 0  # Scale to zero when idle
  maxReplicaCount: 10
  triggers:
  - type: kafka
    metadata:
      topic: report-requests
```

**Savings:** Pay only when reports are being generated

### 4. Cache First
- Put hot reads in Redis
- Use short TTL (30-120s) for balances/history during peaks

**Example:**
```
Without cache:
- DB queries: 100K/min
- DB tier: db.r5.4xlarge ($1,200/month)

With cache (80% hit rate):
- DB queries: 20K/min
- DB tier: db.r5.xlarge ($300/month)
- Redis: $100/month

Savings: $800/month (67% reduction)
```

### 5. Sharding + Replication
- Shard to scale writes
- Replicate inside each shard for reads/HA
- Avoid over-sizing single "monster" DB

### 6. Placement & Affinity
- Keep Payment ↔ Fraud ↔ Wallet in same AZ
- Spread replicas across AZs for HA only where needed

**Example:**
```yaml
affinity:
  podAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchLabels:
            app: payment-api
        topologyKey: topology.kubernetes.io/zone
        # Prefer same AZ as payment-api
```

### 7. Spot for Stateless
- Run background fraud ML scoring on spot instances
- Run async notification workers on spot
- Implement graceful termination

**Example:**
```yaml
nodeSelector:
  node.kubernetes.io/lifecycle: spot  # Use spot instances

# Handle spot termination
lifecycle:
  preStop:
    exec:
      command: ["/bin/sh", "-c", "sleep 30"]  # Drain gracefully
```

**Savings:** 70-90% compute cost for batch workloads

### 8. Reserved Capacity
- Pre-buy the baseline
- Let autoscaling cover bursts

**Example:**
```
Baseline: 5 pods × 24×7 = 3,600 pod-hours/month
Reserved instances: $400/month (40% discount)

Peak burst: 25 pods × 3h × 5 days = 375 pod-hours/month
On-demand: $100/month

Total: $500/month

Versus all on-demand: $900/month
Savings: $400/month (44% reduction)
```

## A Simple Cost Model (Back-of-Napkin)

### Assumptions

**Payment pod needs:**
- 250m CPU
- 256Mi RAM
- Steady state: 5 pods all day
- Peak: 30 pods for 3 hours/day (salary week: 5 days/month)

**Node price:** $0.10 per vCPU-hour (example rate)

### Compute Hours/Month

**Baseline:**
- 5 pods × 24h × 30d = 3,600 pod-hours

**Peak burst:**
- (30 − 5) = 25 extra pods × 3h × 5d = 375 pod-hours

**Total:**
- 3,975 pod-hours × 0.25 vCPU each = 994 vCPU-hours
- 994 × $0.10 = **$99.40/month**

### If You Ran 30 Pods 24×7

- 30 × 24 × 30 = 21,600 pod-hours × 0.25 vCPU = 5,400 vCPU-hours
- 5,400 × $0.10 = **$540/month**

### Elastic vs Always-On

**Savings:** $540 − $99.40 = **$440.60/month (82% reduction)**

**Your numbers will vary, but the shape holds: autoscale beats always-on.**

### Network Egress Sanity Check

**Scenario:** Each payment round-trip moves ~5 KB cross-AZ

**Calculation:**
- 200K transactions/hour × 3 hours = 600K transactions
- 600K × 5 KB = 3 GB cross-AZ traffic
- 3 GB × $0.02/GB = **$0.06 per salary day**

**Multiply by regions, services, and you can see costs add up.**

**Solution:** Same-AZ placement saves this cost entirely.

### DB vs Cache

**Without cache:**
- DB queries: 100K/minute
- DB tier: db.r5.4xlarge = **$1,200/month**

**With cache (80% hit rate):**
- DB queries: 20K/minute
- DB tier: db.r5.xlarge = **$300/month**
- Redis cluster = **$100/month**
- **Total: $400/month**

**Savings: $800/month (67% reduction)**

## Unit Economics: Cost Per 1K Payments

**Executives and product teams understand unit cost, not CPU-hours.**

### Define CPT (Cost Per Thousand Payments)

```
CPT = (Compute + DB + Cache + Network + Observability) / (Payments / 1,000)
```

### Track CPT Weekly

**Example metrics:**
- January Week 1: CPT = $0.05 ✅
- January Week 2: CPT = $0.12 ⚠️ (240% increase!)

**Investigate:**
- Cache hit-rate drop? (90% → 60%)
- Cross-AZ spikes? (chatty service moved)
- Noisy canary? (new version has bug)

### Aim to Reduce CPT

**While holding SLOs:**
- p95 latency < 500ms ✅
- Success rate > 99.9% ✅

**Result:** Improved efficiency without sacrificing quality

## Short, Practical Checklist

### Daily
- [ ] Monitor spend dashboard per service
- [ ] Alert on daily spend > 15% above baseline

### Weekly
- [ ] Review autoscaling (min/max, utilization targets)
- [ ] Check top 5 spenders (compute, DB, network, obs)
- [ ] Review CPT (cost per 1K payments)

### Monthly
- [ ] Right-size top 5 spenders (requests/limits, instance type)
- [ ] Tune observability (sampling, retention)
- [ ] Review reserved capacity vs actual usage
- [ ] Load test and recalibrate SLOs

### Quarterly
- [ ] Full cost audit (tag everything: service, team, env)
- [ ] Optimize database tier and caching strategy
- [ ] Review spot instance candidates
- [ ] Report savings to stakeholders

## Cost Optimization Strategies by Service Type

### Stateless APIs (Payment, Balance Check)
- ✅ HPA with aggressive scale-down
- ✅ Spot instances for non-critical environments
- ✅ Same-AZ placement for chatty services
- ✅ Reserved instances for baseline

### Databases
- ✅ Right-size based on actual usage
- ✅ Read replicas for queries
- ✅ Redis cache to offload reads
- ✅ Archive old data to cold storage

### Message Queues (Kafka, RabbitMQ)
- ✅ Tune retention policies (7 days vs 90 days)
- ✅ Compress messages
- ✅ Right-size broker instances
- ✅ Delete unused topics

### Observability (Logs, Metrics, Traces)
- ✅ Sample traces (1-10% of requests)
- ✅ Log at INFO level in production
- ✅ Downsample metrics over time
- ✅ Set per-service log budgets

### Storage (S3, GCS)
- ✅ Lifecycle policies (hot → warm → cold → delete)
- ✅ Compress files
- ✅ Delete unused buckets
- ✅ Use appropriate storage classes

## Real-World Cost Comparison

### FinPay Wallet: Monthly Cost Breakdown

**Without optimization (over-provisioned):**
```
Compute (30 pods 24×7):        $2,160
Database (db.r5.4xlarge):      $1,200
Redis:                         $100
Network (cross-AZ):            $200
Observability (DEBUG logs):    $800
────────────────────────────────────
TOTAL:                         $4,460/month
```

**With optimization:**
```
Compute (elastic 5-30 pods):   $400   (82% savings)
Database (db.r5.xlarge):       $300   (75% savings)
Redis (with caching):          $100   (same)
Network (same-AZ):             $20    (90% savings)
Observability (INFO, sampled): $150   (81% savings)
────────────────────────────────────
TOTAL:                         $970/month (78% overall savings)
```

**Annual savings: $41,880**

**While maintaining:**
- p95 latency: 320ms (SLO: <500ms) ✅
- Success rate: 99.95% (SLO: >99.9%) ✅
- Peak capacity: 200K transactions/hour ✅

## Key Takeaway

**Cost optimization is not about being cheap—it's about being efficient.**

For FinPay Wallet:
- **Same performance** at 78% lower cost
- **Same reliability** with proper capacity planning
- **Better agility** with elastic scaling

**Three principles:**

1. **Right-size everything** — measure, don't guess
2. **Automate scaling** — grow with demand, shrink when idle
3. **Track unit economics** — cost per transaction, not just total spend

**Cost-aware design enables:**
- Competitive pricing for customers
- Sustainable growth for business
- Focus on features, not firefighting bills

**The best architecture delivers business value efficiently — fast enough, reliable enough, at the right price.**