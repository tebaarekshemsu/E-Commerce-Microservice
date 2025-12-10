BROKER_BINARY=brokerApp

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"
