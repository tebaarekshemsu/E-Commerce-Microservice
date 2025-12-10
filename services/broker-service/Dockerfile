FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o brokerApp ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/brokerApp .

EXPOSE 80

CMD ["./brokerApp"]
