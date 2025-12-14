import { Order, Cart } from '../models/index.js';
import OrderDto from '../dto/OrderDto.js';
import OrderMappingHelper from '../helpers/OrderMappingHelper.js';
import OrderNotFoundException from '../exceptions/OrderNotFoundException.js';
import mongoose from 'mongoose';

class OrderService {
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
    console.log('*** OrderDto, service; save order *');
    const orderData = OrderMappingHelper.mapToEntity(orderDto);
    const order = await Order.create(orderData);
    
    // Reload with populated cart
    const savedOrder = await Order.findById(order._id)
      .populate('cartId', 'cartId userId')
      .exec();
    
    return OrderMappingHelper.mapToDto(savedOrder);
  }

  async update(orderDto) {
    console.log('*** OrderDto, service; update order *');
    if (!orderDto.orderId) {
      throw new OrderNotFoundException('Order ID is required for update');
    }
    
    // First verify the order exists
    await this.findById(orderDto.orderId);
    
    const orderData = OrderMappingHelper.mapToEntity(orderDto);
    delete orderData.orderId; // Remove orderId from update data
    
    const order = await Order.findByIdAndUpdate(
      orderDto.orderId,
      { $set: orderData },
      { new: true, runValidators: true }
    ).populate('cartId', 'cartId userId').exec();
    
    return OrderMappingHelper.mapToDto(order);
  }

  async updateById(orderId, orderDto) {
    console.log('*** OrderDto, service; update order with orderId *');
    // First verify the order exists
    await this.findById(orderId);
    
    const orderData = OrderMappingHelper.mapToEntity(orderDto);
    delete orderData.orderId; // Remove orderId from update data
    
    const order = await Order.findByIdAndUpdate(
      orderId,
      { $set: orderData },
      { new: true, runValidators: true }
    ).populate('cartId', 'cartId userId').exec();
    
    return OrderMappingHelper.mapToDto(order);
  }

  async deleteById(orderId) {
    console.log('*** Void, service; delete order by id *');
    await this.findById(orderId); // Verify order exists
    
    await Order.findByIdAndDelete(orderId);
    
    return true;
  }
}

export default new OrderService();
