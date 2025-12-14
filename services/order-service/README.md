# Order Service

Node.js/Express microservice for managing orders and carts in an e-commerce system.

## Features

- RESTful API for Orders and Carts management
- MongoDB database integration using Mongoose ODM
- External API integration with User Service
- Input validation using express-validator
- Error handling middleware

## Prerequisites

- Node.js 18+ 
- MongoDB 5.0+ or MongoDB 7.0+
- npm or yarn

## Installation

1. Install dependencies:
```bash
npm install
```

2. Create a `.env` file in the root directory:
```env
NODE_ENV=development
PORT=8300
MONGODB_URI=mongodb://localhost:27017/ecommerce_dev_db
DB_HOST=localhost
DB_PORT=27017
DB_NAME=ecommerce_dev_db
DB_USER=
DB_PASSWORD=
USER_SERVICE_URL=http://localhost:8100/user-service/api/users
CONTEXT_PATH=/order-service
```

**Note:** For MongoDB with authentication, use:
```env
MONGODB_URI=mongodb://username:password@localhost:27017/ecommerce_dev_db?authSource=admin
```

## Running the Application

### Development Mode
```bash
npm run dev
```

### Production Mode
```bash
npm start
```

The service will start on port 8300 by default.

## API Endpoints

### Orders

- `GET /order-service/api/orders` - Get all orders
- `GET /order-service/api/orders/:orderId` - Get order by ID
- `POST /order-service/api/orders` - Create a new order
- `PUT /order-service/api/orders` - Update an order
- `PUT /order-service/api/orders/:orderId` - Update an order by ID
- `DELETE /order-service/api/orders/:orderId` - Delete an order

### Carts

- `GET /order-service/api/carts` - Get all carts
- `GET /order-service/api/carts/:cartId` - Get cart by ID
- `POST /order-service/api/carts` - Create a new cart
- `PUT /order-service/api/carts` - Update a cart
- `PUT /order-service/api/carts/:cartId` - Update a cart by ID
- `DELETE /order-service/api/carts/:cartId` - Delete a cart

## Docker

Build the Docker image:
```bash
docker build -t order-service:0.1.0 .
```

Run with Docker Compose (includes MongoDB):
```bash
docker-compose up
```

This will start both the order service and MongoDB in containers.

## Project Structure

```
order-service/
├── src/
│   ├── config/         # Application configuration
│   ├── controllers/    # Route controllers
│   ├── dto/           # Data Transfer Objects
│   ├── exceptions/    # Custom exceptions
│   ├── helpers/       # Mapping helpers
│   ├── middleware/    # Express middleware
│   ├── models/        # Mongoose models
│   ├── routes/        # Express routes
│   └── services/      # Business logic
├── server.js          # Application entry point
└── package.json       # Dependencies
```

## Environment Variables

- `NODE_ENV` - Environment (development/production)
- `PORT` - Server port (default: 8300)
- `MONGODB_URI` - MongoDB connection string (overrides other DB settings)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 27017)
- `DB_NAME` - Database name (default: ecommerce_dev_db)
- `DB_USER` - Database user (optional, for authentication)
- `DB_PASSWORD` - Database password (optional, for authentication)
- `USER_SERVICE_URL` - User service API URL
- `CONTEXT_PATH` - API context path (default: /order-service)

## Database Schema

### Cart Collection
```javascript
{
  _id: ObjectId,
  userId: Number,
  createdAt: Date,
  updatedAt: Date
}
```

### Order Collection
```javascript
{
  _id: ObjectId,
  orderDate: Date,
  orderDesc: String,
  orderFee: Number,
  cartId: ObjectId (reference to Cart),
  createdAt: Date,
  updatedAt: Date
}
```
