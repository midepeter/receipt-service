# === Build Stage ===
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o receipt-service

# === Runtime Stage ===
FROM debian:bullseye-slim

# Install wkhtmltopdf and clean up
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/receipt-service .
COPY --from=builder /app/templates/receiptnewtwo.html ./templates/receiptnewtwo.html
COPY --from=builder /app/templates/transaction-history-new.html ./templates/transaction-history-new.html

EXPOSE 4002

ENTRYPOINT ["./receipt-service"]
