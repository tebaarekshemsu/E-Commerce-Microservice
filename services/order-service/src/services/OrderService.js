import { Order, Cart } from '../models/index.js';
import OrderDto from '../dto/OrderDto.js';
import OrderMappingHelper from '../helpers/OrderMappingHelper.js';
import OrderNotFoundException from '../exceptions/OrderNotFoundException.js';
import grpcClients from './grpcClients.js';
import rabbitmqService from './rabbitmq.js';
import mongoose from 'mongoose';

class OrderService {
  constructor() {
    // Initialize RabbitMQ connection
    this.initializeRabbitMQ();
  }

  async initializeRabbitMQ() {
    try {
      await rabbitmqService.connect();
    } catch (error) {
      console.error('Failed to initialize RabbitMQ:', error);
    }
  }

  async findAll() {
    console.log('*** OrderDto List, service; fetch all orders *');
    const orders = await Order.find()
      .populate('cartId', 'cartId userId')
      .exec();
    
    return orders
      .map(order => OrderMappingHelper.mapToDto(order))
      .filter(order => order !== null);
  }

  async findById(orderId) {
    console.log('*** OrderDto, service; fetch order by id *');
    
    if (!mongoose.Types.ObjectId.isValid(orderId)) {
      throw new OrderNotFoundException(`Order with id: ${orderId} not found`);
    }
    
    const order = await Order.findById(orderId)
      .populate('cartId', 'cartId userId')
      .exec();
    
    if (!order) {
      throw new OrderNotFoundException(`Order with id: ${orderId} not found`);
    }
    
    return OrderMappingHelper.mapToDto(order);
  }

  async save(orderDto) {
    console.log('*** OrderDto, service; save order with gRPC and RabbitMQ integration *');
    
    try {
      // Step 1: Validate cart if provided
      if (orderDto.cartId) {
        const cart = await Cart.findById(orderDto.cartId);
        if (!cart) {
          throw new Error(`Cart with id ${orderDto.cartId} not found`);
        }
      }

      // Step 2: If we have product information, check availability via gRPC
      if (orderDto.productId && orderDto.quantity) {
        try {
          const availabilityCheck = await grpcClients.checkProductAvailability(
            orderDto.productId, 
            orderDto.quantity
          );
          
          if (!availabilityCheck.available) {
            throw new Error(`Product ${orderDto.productId} is not available in requested quantity`);
          }

          // Get product details for pricing
          const product = await grpcClients.getProduct(orderDto.productId);
          if (product && !orderDto.orderFee) {
            orderDto.orderFee = parseFloat(product.price) * orderDto.quantity;
          }
        } catch (grpcError) {
          console.warn('gRPC product check failed, proceeding without validation:', grpcError.message);
        }
      }

      // Step 3: Create the order
      const orderData = OrderMappingHelper.mapToEntity(orderDto);
      const order = await Order.create(orderData);
      
      // Reload with populated cart
      const savedOrder = await Order.findById(order._id)
        .populate('cartId', 'cartId userId')
        .exec();

      const orderDtoResult = OrderMappingHelper.mapToDto(savedOrder);

      // Step 4: Process payment via gRPC (if order fee exists)
      if (orderDtoResult.orderFee && orderDtoResult.orderFee > 0) {
        try {
          const paymentResult = await grpcClients.processPayment({
            orderId: orderDtoResult.orderId,
            orderFee: orderDtoResult.orderFee,
            currency: 'USD',
            paymentMethod: 'credit_card'
          });

          console.log('Payment processing result:', paymentResult);
          
          if (!paymentResult.success) {
            // Update order status or handle payment failure
            console.warn('Payment failed for order:', orderDtoResult.orderId);
          }
        } catch (paymentError) {
          console.error('Payment processing failed:', paymentError.message);
        }
      }

      // Step 5: Publish order created event to RabbitMQ
      try {
        await rabbitmqService.publishOrderEvent('order_created', {
          orderId: orderDtoResult.orderId,
          userId: orderDtoResult.cartId?.userId || 'anonymous',
          email: 'customer@example.com', // In real app, get from user service
          orderFee: orderDtoResult.orderFee,
          items: orderDto.items || [],
          timestamp: new Date().toISOString()
        });
      } catch (eventError) {
        console.error('Failed to publish order event:', eventError.message);
      }

      return orderDtoResult;
    } catch (error) {
      console.error('Error in order creation:', error);
      throw error;
    }
  }

  async update(orderDto) {
    console.log('*** OrderDto, service; update order *');
    if (!orderDto.orderId) {
      throw new OrderNotFoundException('Order ID is required for update');
    }
    
    // First verify the order exists
    const existingOrder = await this.findById(orderDto.orderId);
    
    const orderData = OrderMappingHelper.mapToEntity(orderDto);
    delete orderData.orderId; // Remove orderId from update data
    
    const order = await Order.findByIdAndUpdate(
      orderDto.orderId,
      { $set: orderData },
      { new: true, runValidators: true }
    ).populate('cartId', 'cartId userId').exec();
    
    const updatedOrderDto = OrderMappingHelper.mapToDto(order);

    // Publish order updated event
    try {
      await rabbitmqService.publishOrderEvent('order_updated', {
        orderId: updatedOrderDto.orderId,
        userId: updatedOrderDto.cartId?.userId || 'anonymous',
        orderFee: updatedOrderDto.orderFee,
        previousState: existingOrder,
        timestamp: new Date().toISOString()
      });
    } catch (eventError) {
      console.error('Failed to publish order update event:', eventError.message);
    }
    
    return updatedOrderDto;
  }

  async updateById(orderId, orderDto) {
    console.log('*** OrderDto, service; update order with orderId *');
    // First verify the order exists
    const existingOrder = await this.findById(orderId);
    
    const orderData = OrderMappingHelper.mapToEntity(orderDto);
    delete orderData.orderId; // Remove orderId from update data
    
    const order = await Order.findByIdAndUpdate(
      orderId,
      { $set: orderData },
      { new: true, runValidators: true }
    ).populate('cartId', 'cartId userId').exec();
    
    const updatedOrderDto = OrderMappingHelper.mapToDto(order);

    // Publish order updated event
    try {
      await rabbitmqService.publishOrderEvent('order_updated', {
        orderId: updatedOrderDto.orderId,
        userId: updatedOrderDto.cartId?.userId || 'anonymous',
        orderFee: updatedOrderDto.orderFee,
        previousState: existingOrder,
        timestamp: new Date().toISOString()
      });
    } catch (eventError) {
      console.error('Failed to publish order update event:', eventError.message);
    }
    
    return updatedOrderDto;
  }

  async deleteById(orderId) {
    console.log('*** Void, service; delete order by id *');
    const existingOrder = await this.findById(orderId); // Verify order exists
    
    await Order.findByIdAndDelete(orderId);

    // Publish order cancelled event
    try {
      await rabbitmqService.publishOrderEvent('order_cancelled', {
        orderId: existingOrder.orderId,
        userId: existingOrder.cartId?.userId || 'anonymous',
        orderFee: existingOrder.orderFee,
        timestamp: new Date().toISOString()
      });
    } catch (eventError) {
      console.error('Failed to publish order cancellation event:', eventError.message);
    }
    
    return true;
  }

  // New method for order status updates (shipping, delivery, etc.)
  async updateOrderStatus(orderId, status) {
    console.log(`*** Update order ${orderId} status to ${status} *`);
    
    const existingOrder = await this.findById(orderId);
    
    // Update order with new status (you might want to add a status field to your Order model)
    const order = await Order.findByIdAndUpdate(
      orderId,
      { $set: { status: status, updatedAt: new Date() } },
      { new: true, runValidators: true }
    ).populate('cartId', 'cartId userId').exec();

    const updatedOrderDto = OrderMappingHelper.mapToDto(order);

    // Publish appropriate event based on status
    let eventType = 'order_updated';
    if (status === 'shipped') eventType = 'order_shipped';
    else if (status === 'delivered') eventType = 'order_delivered';
    else if (status === 'cancelled') eventType = 'order_cancelled';

    try {
      await rabbitmqService.publishOrderEvent(eventType, {
        orderId: updatedOrderDto.orderId,
        userId: updatedOrderDto.cartId?.userId || 'anonymous',
        orderFee: updatedOrderDto.orderFee,
        status: status,
        timestamp: new Date().toISOString()
      });
    } catch (eventError) {
      console.error('Failed to publish order status event:', eventError.message);
    }

    return updatedOrderDto;
  }
}

export default new OrderService();
