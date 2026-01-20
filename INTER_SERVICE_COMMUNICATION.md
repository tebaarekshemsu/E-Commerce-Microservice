# Inter-Service Communication Implementation

This document describes the complete inter-service communication implementation using gRPC for synchronous calls and RabbitMQ for asynchronous events.

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Order Service │    │ Product Service │    │ Payment Service │
│   (Node.js)     │    │     (Go)        │    │     (Go)        │
│                 │    │                 │    │                 │
│ HTTP: 8083      │    │ HTTP: 8082      │    │ HTTP: 8085      │
│ gRPC Client     │    │ gRPC: 9001      │    │ gRPC: 9002      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │  User Service   │
                    │   (Python)      │
                    │                 │
                    │ HTTP: 8001      │
                    │ gRPC: 9003      │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Notification    │    │    RabbitMQ     │    │  Broker Service │
│   Service       │    │   (Message      │    │   (API Gateway) │
│    (Go)         │    │    Broker)      │    │     (Go)        │
│                 │    │                 │    │                 │
│ HTTP: 8084      │    │ AMQP: 5672      │    │ HTTP: 80        │
│ RabbitMQ        │    │ Mgmt: 15672     │    │ Reverse Proxy   │
│ Consumer        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🔄 Communication Patterns

### 1. Synchronous Communication (gRPC)

#### Product Service gRPC Server (Port 9001)
- **GetProduct**: Retrieve product details by ID
- **CheckAvailability**: Verify stock availability for quantity

#### Payment Service gRPC Server (Port 9002)  
- **ProcessPayment**: Process payment transactions

#### User Service gRPC Server (Port 9003)
- **GetUser**: Retrieve user profile information

### 2. Asynchronous Communication (RabbitMQ)

#### Event Queues:
- **orders**: Order lifecycle events (created, updated, shipped, delivered, cancelled)
- **users**: User events (registered, password_reset_requested)
- **inventory**: Inventory events (low_stock alerts)

#### Event Publishers:
- **Order Service**: Publishes order events
- **User Service**: Publishes user events (future implementation)
- **Product Service**: Publishes inventory events (future implementation)

#### Event Consumers:
- **Notification Service**: Consumes all events and sends notifications

## 📋 Implementation Details

### Order Service (Node.js)

#### gRPC Clients (`src/services/grpcClients.js`)
```javascript
// Product availability check
const availability = await grpcClients.checkProductAvailability(productId, quantity);

// Payment processing
const paymentResult = await grpcClients.processPayment({
  orderId: order.id,
  amount: order.total,
  currency: 'USD'
});

// User information
const user = await grpcClients.getUser(userId);
```

#### RabbitMQ Publisher (`src/services/rabbitmq.js`)
```javascript
// Publish order events
await rabbitmqService.publishOrderEvent('order_created', {
  orderId: order.id,
  userId: order.userId,
  email: user.email,
  total: order.total
});
```

### Product Service (Go)

#### gRPC Server (`grpc/server.go`)
```go
func (s *ProductServer) CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error) {
    // Check product availability in database
    product, err := s.models.Product.GetOne(productId)
    available := product.Quantity >= int(req.Quantity)
    
    return &pb.CheckAvailabilityResponse{Available: available}, nil
}
```

### Payment Service (Go)

#### gRPC Server (`grpc/server.go`)
```go
func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
    // Create payment record
    payment := data.Payment{
        OrderID: orderID,
        IsPayed: false,
        PaymentStatus: "processing",
    }
    
    // Process payment logic
    success := processPaymentLogic(req)
    
    return &pb.PaymentResponse{
        Success: success,
        TransactionId: generateTransactionId(),
        Message: "Payment processed",
    }, nil
}
```

### Notification Service (Go)

#### RabbitMQ Consumer (`internal/handlers/notification.go`)
```go
func (h *NotificationHandler) HandleOrderEvent(data []byte) error {
    var event models.OrderEvent
    json.Unmarshal(data, &event)
    
    switch event.EventType {
    case "order_created":
        return h.sendOrderConfirmation(event)
    case "order_shipped":
        return h.sendShippingNotification(event)
    // ... other event types
    }
}
```

## 🚀 Service Startup Sequence

1. **Infrastructure Services**
   - PostgreSQL (port 5432)
   - MongoDB (port 27017)  
   - RabbitMQ (port 5672, management 15672)

2. **Core Services**
   - Product Service (HTTP: 8082, gRPC: 9001)
   - Payment Service (HTTP: 8085, gRPC: 9002)
   - User Service (HTTP: 8001, gRPC: 9003)

3. **Dependent Services**
   - Order Service (HTTP: 8083) - depends on Product, Payment, User services
   - Notification Service (HTTP: 8084) - depends on RabbitMQ

4. **API Gateway**
   - Broker Service (HTTP: 80) - reverse proxy to all services

## 🔧 Configuration

### Environment Variables

#### Order Service
```env
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
PRODUCT_SERVICE_GRPC_URL=product-service:9001
PAYMENT_SERVICE_GRPC_URL=payment-service:9002
USER_SERVICE_GRPC_URL=user-service:9003
```

#### Notification Service
```env
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
```

### Docker Compose Ports
```yaml
services:
  product-service:
    ports:
      - "8082:80"    # HTTP API
      - "9001:9001"  # gRPC Server
      
  payment-service:
    ports:
      - "8085:80"    # HTTP API
      - "9002:9002"  # gRPC Server
      
  user-service:
    ports:
      - "8001:8000"  # HTTP API
      - "9003:9003"  # gRPC Server
      
  order-service:
    ports:
      - "8083:80"    # HTTP API (gRPC client only)
      
  notification-service:
    ports:
      - "8084:8080"  # HTTP API (RabbitMQ consumer)
      
  rabbitmq:
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # Management UI
```

## 📊 Event Flow Examples

### Order Creation Flow
```
1. Client → POST /order-service/api/orders
2. Order Service → gRPC CheckAvailability → Product Service
3. Order Service → Create Order in MongoDB
4. Order Service → gRPC ProcessPayment → Payment Service
5. Order Service → Publish "order_created" → RabbitMQ
6. Notification Service → Consume event → Send confirmation email
7. Order Service → Return order response to client
```

### Order Status Update Flow
```
1. Admin → PUT /order-service/api/orders/{id}/status
2. Order Service → Update order status
3. Order Service → Publish "order_shipped" → RabbitMQ
4. Notification Service → Consume event → Send shipping notification
```

## 🧪 Testing Communication

### Health Check Endpoints
- Order Service: `GET /health` - Shows gRPC and RabbitMQ connection status
- Individual Services: `GET /ping` - Basic health check

### Manual Testing
```bash
# Test gRPC directly (requires grpcurl)
grpcurl -plaintext -d '{"id":"1"}' localhost:9001 product.ProductService/GetProduct

# Test RabbitMQ Management UI
http://localhost:15672 (guest/guest)

# Test order creation with communication
curl -X POST http://localhost/order-service/api/orders \
  -H "Content-Type: application/json" \
  -d '{"orderDesc":"Test Order","orderFee":99.99,"productId":"1","quantity":2}'
```

## 🔍 Monitoring & Debugging

### Logs to Monitor
- gRPC call logs in Product/Payment/User services
- RabbitMQ message publishing in Order Service
- RabbitMQ message consumption in Notification Service
- Connection status logs for all services

### Common Issues
1. **gRPC Connection Refused**: Check if target service is running and port is exposed
2. **RabbitMQ Connection Failed**: Verify RabbitMQ is running and URL is correct
3. **Proto Compilation Errors**: Ensure proto files are properly generated
4. **Message Not Consumed**: Check queue names and RabbitMQ consumer status

## 🚀 Deployment

### Development
```bash
cd project
docker-compose up --build
```

### Production Considerations
- Use proper gRPC load balancing
- Implement circuit breakers for gRPC calls
- Add message persistence and dead letter queues for RabbitMQ
- Monitor gRPC and RabbitMQ metrics
- Implement proper authentication for gRPC services

## 📈 Future Enhancements

1. **Service Mesh**: Consider Istio for advanced traffic management
2. **Event Sourcing**: Implement event store for audit trails
3. **CQRS**: Separate read/write models for better scalability
4. **Distributed Tracing**: Add OpenTelemetry for request tracing
5. **Schema Registry**: Implement Avro/Protobuf schema evolution