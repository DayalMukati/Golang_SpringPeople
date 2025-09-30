# Strangler Fig Pattern Example

This example demonstrates the Strangler Fig Pattern for incrementally modernizing a legacy monolith.

## Services

1. **strangler-proxy** (front door) - Runs on port 8080
   - Single entry point for all client requests
   - Routes /api/users/* → new users service
   - Routes everything else → legacy app
   - Adds request ID for tracing

2. **usersvc** (new microservice) - Runs on port 7001
   - Handles only /api/users/* endpoints
   - Represents the first migrated slice
   - Uses new datastore (simulated)

3. **legacy** (monolith mock) - Runs on port 7000
   - Handles everything not yet migrated
   - Still serves /api/payments/*, /api/refunds/*, etc.

## How to Run

1. Start legacy:
```bash
cd legacy
go run main.go
```

2. Start usersvc:
```bash
cd usersvc
go run main.go
```

3. Start strangler-proxy:
```bash
cd strangler-proxy
go run main.go
```

4. Test routing:
```bash
# Goes to NEW users service
curl -s http://localhost:8080/api/users/123

# Goes to LEGACY monolith (not migrated yet)
curl -s http://localhost:8080/api/payments/tx/abc
```

## Key Points

- One stable URL for clients (localhost:8080)
- Different backends per feature
- Incremental migration (one feature at a time)
- Low risk - roll back individual routes if needed