# === Build Stage ===
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o receipt-service

# === Runtime Stage ===
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/receipt-service .
EXPOSE 4000
ENTRYPOINT ["./receipt-service"]
