import dotenv from 'dotenv';
dotenv.config();
import express from 'express';
import cors from 'cors';
import morgan from 'morgan';
import { connectDB } from './src/config/database.js';

// Load models
import './src/models/index.js';

import orderRoutes from './src/routes/orderRoutes.js';
import cartRoutes from './src/routes/cartRoutes.js';
import errorHandler from './src/middleware/errorHandler.js';

// Import services for initialization
import rabbitmqService from './src/services/rabbitmq.js';
import grpcClients from './src/services/grpcClients.js';

const app = express();
const PORT = process.env.PORT || 8300;
const CONTEXT_PATH = process.env.CONTEXT_PATH || '/order-service';

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(morgan('dev'));

// Health check endpoint
app.get('/', (req, res) => {
  res.json({ message: 'Order controller responding!!' });
});

// Health check with service status
app.get('/health', async (req, res) => {
  try {
    const grpcHealth = await grpcClients.healthCheck();
    const rabbitmqHealth = rabbitmqService.isConnected;
    
    res.json({
      status: 'OK',
      timestamp: new Date().toISOString(),
      services: {
        database: 'connected',
        rabbitmq: rabbitmqHealth ? 'connected' : 'disconnected',
        grpc: grpcHealth
      }
    });
  } catch (error) {
    res.status(503).json({
      status: 'ERROR',
      timestamp: new Date().toISOString(),
      error: error.message
    });
  }
});

// Routes
app.use(`${CONTEXT_PATH}/api/orders`, orderRoutes);
app.use(`${CONTEXT_PATH}/api/carts`, cartRoutes);

// Error handling middleware
app.use(errorHandler);

// Database connection and server start
const startServer = async () => {
  try {
    await connectDB();
    
    // Initialize RabbitMQ connection
    try {
      await rabbitmqService.connect();
      console.log('✅ RabbitMQ connected successfully');
    } catch (rabbitmqError) {
      console.warn('⚠️ RabbitMQ connection failed:', rabbitmqError.message);
    }
    
    app.listen(PORT, () => {
      console.log(`🚀 Order Service is running on port ${PORT}`);
      console.log(`📍 Context path: ${CONTEXT_PATH}`);
      console.log(`🔗 gRPC clients initialized for Product, Payment, and User services`);
    });
  } catch (error) {
    console.error('Unable to start the server:', error);
    process.exit(1);
  }
};

// Graceful shutdown
process.on('SIGINT', async () => {
  console.log('Shutting down Order Service...');
  await rabbitmqService.close();
  process.exit(0);
});

process.on('SIGTERM', async () => {
  console.log('Shutting down Order Service...');
  await rabbitmqService.close();
  process.exit(0);
});

startServer();

export default app;
