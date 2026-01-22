FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY migrations ./migrations
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o order-service ./cmd/order-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./order-service"]
