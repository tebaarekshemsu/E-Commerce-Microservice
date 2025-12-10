FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o paymentApp ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/paymentApp .

EXPOSE 80

CMD ["./paymentApp"]
