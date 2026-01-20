#!/bin/sh
set -e

echo "Installing dependencies..."
# We need python3 and pip, plus the grpc tools
apk add --no-cache python3 py3-pip protobuf-dev
# Create a venv to avoid PEP 668 restrictions
python3 -m venv /tmp/venv
. /tmp/venv/bin/activate

echo "Installing Python gRPC tools..."
pip install grpcio-tools

echo "Generating User Service protos..."
# The output directory must exist
mkdir -p services/user-service/proto

# python -m grpc_tools.protoc 
#   -I project/protos          (Includes path)
#   --python_out=...           (Output for messages)
#   --grpc_python_out=...      (Output for services)
#   project/protos/user.proto  (Input file)

python -m grpc_tools.protoc \
    -I project/protos \
    --python_out=services/user-service/proto \
    --grpc_python_out=services/user-service/proto \
    project/protos/user.proto

echo "Done!"
