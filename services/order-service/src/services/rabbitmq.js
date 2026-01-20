import amqp from 'amqplib';

class RabbitMQService {
  constructor() {
    this.connection = null;
    this.channel = null;
    this.isConnected = false;
  }

  async connect() {
    try {
      const rabbitmqUrl = process.env.RABBITMQ_URL || 'amqp://guest:guest@localhost:5672/';
      console.log('🔌 Connecting to RabbitMQ:', rabbitmqUrl);
      
      this.connection = await amqp.connect(rabbitmqUrl);
      this.channel = await this.connection.createChannel();
      this.isConnected = true;
      
      console.log('✅ Connected to RabbitMQ');
      
      // Handle connection errors
      this.connection.on('error', (err) => {
        console.error('RabbitMQ connection error:', err);
        this.isConnected = false;
      });
      
      this.connection.on('close', () => {
        console.log('RabbitMQ connection closed');
        this.isConnected = false;
      });
      
    } catch (error) {
      console.error('Failed to connect to RabbitMQ:', error);
      this.isConnected = false;
      throw error;
    }
  }

  async ensureQueue(queueName) {
    if (!this.isConnected || !this.channel) {
      await this.connect();
    }
    
    await this.channel.assertQueue(queueName, {
      durable: true
    });
  }

  async publishOrderEvent(eventType, orderData) {
    try {
      const queueName = 'orders';
      await this.ensureQueue(queueName);
      
      const event = {
        eventType,
        orderId: orderData.orderId,
        userId: orderData.userId,
        email: orderData.email || 'customer@example.com', // Default email for demo
        items: orderData.items || [],
        total: orderData.orderFee || 0,
        timestamp: new Date().toISOString(),
        metadata: orderData
      };
      
      const message = Buffer.from(JSON.stringify(event));
      
      const published = this.channel.sendToQueue(queueName, message, {
        persistent: true
      });
      
      if (published) {
        console.log(`📤 Published ${eventType} event for order ${orderData.orderId}`);
      } else {
        console.warn(`⚠️ Failed to publish ${eventType} event for order ${orderData.orderId}`);
      }
      
      return published;
    } catch (error) {
      console.error('Error publishing order event:', error);
      throw error;
    }
  }

  async publishUserEvent(eventType, userData) {
    try {
      const queueName = 'users';
      await this.ensureQueue(queueName);
      
      const event = {
        eventType,
        userId: userData.userId,
        email: userData.email,
        name: userData.name || 'Customer',
        timestamp: new Date().toISOString(),
        metadata: userData
      };
      
      const message = Buffer.from(JSON.stringify(event));
      
      const published = this.channel.sendToQueue(queueName, message, {
        persistent: true
      });
      
      if (published) {
        console.log(`📤 Published ${eventType} event for user ${userData.userId}`);
      }
      
      return published;
    } catch (error) {
      console.error('Error publishing user event:', error);
      throw error;
    }
  }

  async publishInventoryEvent(eventType, inventoryData) {
    try {
      const queueName = 'inventory';
      await this.ensureQueue(queueName);
      
      const event = {
        eventType,
        productId: inventoryData.productId,
        productName: inventoryData.productName,
        quantity: inventoryData.quantity,
        threshold: inventoryData.threshold || 10,
        timestamp: new Date().toISOString(),
        metadata: inventoryData
      };
      
      const message = Buffer.from(JSON.stringify(event));
      
      const published = this.channel.sendToQueue(queueName, message, {
        persistent: true
      });
      
      if (published) {
        console.log(`📤 Published ${eventType} event for product ${inventoryData.productId}`);
      }
      
      return published;
    } catch (error) {
      console.error('Error publishing inventory event:', error);
      throw error;
    }
  }

  async close() {
    try {
      if (this.channel) {
        await this.channel.close();
      }
      if (this.connection) {
        await this.connection.close();
      }
      this.isConnected = false;
      console.log('RabbitMQ connection closed');
    } catch (error) {
      console.error('Error closing RabbitMQ connection:', error);
    }
  }
}

// Create singleton instance
const rabbitmqService = new RabbitMQService();

export default rabbitmqService;