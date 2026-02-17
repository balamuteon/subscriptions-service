FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/subscriptions-api ./cmd/api

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /bin/subscriptions-api /app/subscriptions-api
COPY --from=builder /app/.env /app/.env

EXPOSE 8080
CMD ["/app/subscriptions-api"]
