# Step 1: Build stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o payment-service .

# Step 2: Runtime stage
FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/payment-service .
EXPOSE 8080
CMD ["./payment-service"]