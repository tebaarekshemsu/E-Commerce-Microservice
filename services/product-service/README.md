# Product Service (Go)

This is a port of the Java Product Service to Go.

## Features

- REST API for Products and Categories
- PostgreSQL Database connection
- Docker support

## Structure

- `cmd/api`: Entry point and HTTP handlers
- `data`: Database models and logic
- `Dockerfile`: Docker build configuration

## Configuration

The service uses the following environment variables:

- `DSN`: Database connection string (e.g., `host=postgres port=5432 user=postgres password=password dbname=product_service sslmode=disable timezone=UTC connect_timeout=5`)
- `PORT`: Web server port (default: 80)

## Running

### Local

1. Ensure PostgreSQL is running and accessible.
2. Set `DSN` environment variable.
3. Run `go run ./cmd/api`

### Docker

1. Build the image: `docker build -t product-service-go .`
2. Run the container: `docker run -p 80:80 -e DSN=... product-service-go`

## API Endpoints

Base path: `/product-service`

### Products

- `GET /api/products`: Get all products
- `GET /api/products/{productId}`: Get product by ID
- `POST /api/products`: Create a new product
- `PUT /api/products`: Update a product
- `PUT /api/products/{productId}`: Update a product by ID
- `DELETE /api/products/{productId}`: Delete a product

### Categories

- `GET /api/categories`: Get all categories
- `GET /api/categories/{categoryId}`: Get category by ID
- `POST /api/categories`: Create a new category
- `PUT /api/categories`: Update a category
- `PUT /api/categories/{categoryId}`: Update a category by ID
- `DELETE /api/categories/{categoryId}`: Delete a category

## Notes

- This service mimics the Java `context-path` of `/product-service`.
- Zipkin tracing and Config Server are not currently implemented in this Go version.
