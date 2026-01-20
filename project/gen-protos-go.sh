#!/bin/sh
set -e

echo "Installing dependencies..."
apk add --no-cache protobuf-dev git > /dev/null

echo "Setting up Go environment..."
export GOPROXY=https://proxy.golang.org,direct
export PATH="$PATH:$(go env GOPATH)/bin"

echo "Installing plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

echo "Generating Product Service protos..."
protoc --proto_path=project/protos \
       --go_out=services/product-service/proto --go_opt=paths=source_relative \
       --go-grpc_out=services/product-service/proto --go-grpc_opt=paths=source_relative \
       project/protos/product.proto

echo "Generating Payment Service protos..."
protoc --proto_path=project/protos \
       --go_out=services/payment-service/proto --go_opt=paths=source_relative \
       --go-grpc_out=services/payment-service/proto --go-grpc_opt=paths=source_relative \
       project/protos/payment.proto

echo "Generating Broker Service protos..."
protoc --proto_path=project/protos \
       --go_out=services/broker-service/proto --go_opt=paths=source_relative \
       --go-grpc_out=services/broker-service/proto --go-grpc_opt=paths=source_relative \
       project/protos/product.proto project/protos/payment.proto

echo "Done!"
