# E-Commerce-Microservice

A distributed e-commerce platform based on microservice architecture, supporting modular development, scalability, and event-driven workflows.

## 🚀 Overview

This project implements an event-based, distributed e-commerce solution with a broker design, using a **database-per-service** approach and focusing on loose coupling, resiliency, and scalability. The architecture supports both synchronous (HTTP-based) and asynchronous (message-broker) flows for different service interactions.

- **Homepage:** [E-Commerce-Microservice Vercel App](https://e-commerce-microservice.vercel.app)

## 🏗️ Architecture

- **Microservices:** Independent services (Order, Payment, Broker, etc.)
- **API Gateway:** Nginx routes traffic and manages load balancing (planned/optional).
- **Database per Service:** Each core component (e.g., Order, Product) manages its own PostgreSQL database.
- **Broker:** Facilitates event-driven communication with published topics and standardized event formats.
- **Frontend/Backend Stack:** Initial design referenced Django (backend), Angular (frontend), but the codebase also includes Go and Node.js/Express for some services.
- **Swagger/OpenAPI:** Integrated for API specification and interactive use (docs folder).

## 📦 Main Features & Structure

- **Broker Service:** Written in Go, built using multi-stage Dockerfile (`services/broker-service/Dockerfile`). Exposes HTTP endpoint(s) for event publishing and relaying.
- **Order Service:** Node.js with Express (`services/order-service/server.js`), handles order creation, cart management, and order lifecycles. Includes health checks, error handling middleware, and REST API structure.
- **Payment Service:** Go application managing payments, exposing standard CRUD operations for payment data, and operating on its own DB (`services/payment-service/data/models.go`).
- **Documentation:** Swagger UI served in the `docs/` folder with OpenAPI specs, static assets, and a basic HTML/CSS layout.
- **Docker-Oriented:** `project/Makefile` and Dockerfiles enable orchestration and one-command management of services.
- **Event Flows:** Hybrid synchronous (inventory check) and asynchronous (payment via broker) event processing.
- **API Contracts:** Well-defined OpenAPI specs for clear communication, error, and success response formats.

## 🗂️ Directory Structure (Partial View)

```
docs/
  index.html            # API Documentation (Swagger UI)
  index.css
  swagger-initializer.js
  oauth2-redirect.html
  ...
services/
  broker-service/       # Go-based broker
  order-service/        # Node.js/Express order handling
  payment-service/      # Go-based payment processing
project/
  Makefile              # Docker orchestration commands
README.md
```

*Note: This is a partial directory list. For the full structure, see the [repository on GitHub](https://github.com/tebaarekshemsu/E-Commerce-Microservice).*

## 📝 Project Status & Documentation

- **Phase 1:** System design, component diagrams, and deployment architecture completed ([progress report](docs/weekly-progress-week7-8.md)).
- **API Docs:** See the `docs/` directory for live Swagger UI.
- **Risks:** Docker networking and eventual consistency in asynchronous flows highlighted as current challenges.

## 🛠️ Getting Started

### Prerequisites

- Docker and Docker Compose

### Quick Start

```sh
cd project
make up_build  # Builds and starts all services
make logs      # Streams logs for all services
make down      # Stops all containers
```

### Order Service Dev Start

```sh
cd services/order-service
npm install
node server.js
```

(Broker and payment services should be built and run via Docker.)

## 🧩 Contributing

PRs and suggestions are welcome! Please raise an issue or pull request on GitHub.

## 🪪 License

Currently undefined.

---

For more examples, see the [API Documentation](https://github.com/tebaarekshemsu/E-Commerce-Microservice/tree/main/docs) or the [repository homepage](https://github.com/tebaarekshemsu/E-Commerce-Microservice).
