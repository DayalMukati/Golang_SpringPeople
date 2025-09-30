# Stage 1: Build
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go mod init payment && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o payment .

# Stage 2: Run
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/payment /payment
EXPOSE 8080
ENTRYPOINT ["/payment"]