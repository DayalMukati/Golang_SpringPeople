FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o payment-service .

FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/payment-service .
EXPOSE 8080
CMD ["./payment-service"]