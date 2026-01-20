## up: starts all containers in the background
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## logs: view logs for all services
logs:
	docker-compose logs -f

## proto: generates gRPC code from proto files
proto:
	@echo "Generating gRPC code..."
	docker run --rm -v $(CURDIR):/workspace -w /workspace simplysecurity/protoc \
		--proto_path=protos \
		--go_out=../services/product-service/proto --go_opt=paths=source_relative \
		--go-grpc_out=../services/product-service/proto --go-grpc_opt=paths=source_relative \
		protos/product.proto
	docker run --rm -v $(CURDIR):/workspace -w /workspace simplysecurity/protoc \
		--proto_path=protos \
		--go_out=../services/payment-service/proto --go_opt=paths=source_relative \
		--go-grpc_out=../services/payment-service/proto --go-grpc_opt=paths=source_relative \
		protos/payment.proto
	docker run --rm -v $(CURDIR):/workspace -w /workspace simplysecurity/protoc \
		--proto_path=protos \
		--go_out=../services/broker-service/proto --go_opt=paths=source_relative \
		--go-grpc_out=../services/broker-service/proto --go-grpc_opt=paths=source_relative \
		protos/product.proto protos/payment.proto
	docker run --rm -v $(CURDIR):/workspace -w /workspace simplysecurity/protoc \
		--proto_path=protos \
		--python_out=../services/user-service/proto \
		--grpc_python_out=../services/user-service/proto \
		protos/user.proto
	@echo "gRPC code generated!"
