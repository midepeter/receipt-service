# === Build Stage ===
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o receipt-service

# === Runtime Stage ===
FROM ubuntu:20.04

RUN apt-get update && \
    apt-get install -y wkhtmltopdf

WORKDIR /app
COPY --from=builder /app/receipt-service .
COPY --from=builder /app/templates/receiptnewtwo.html .
COPY --from=builder /app/templates/transaction-history-new.html .

EXPOSE 4002
ENTRYPOINT ["./receipt-service"]
