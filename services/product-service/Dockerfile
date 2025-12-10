FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./ 
# go mod download will fail without network, but assuming it works in build environment
# RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o productApp ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/productApp .

EXPOSE 80

CMD ["./productApp"]
