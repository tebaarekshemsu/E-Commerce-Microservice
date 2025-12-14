import OrderDto from '../dto/OrderDto.js';
import CartDto from '../dto/CartDto.js';
import mongoose from 'mongoose';

class OrderMappingHelper {
  static mapToDto(order) {
    if (!order) return null;
    
    // Convert Mongoose document to plain object
    const orderData = order.toObject ? order.toObject() : order;
    
    // Handle cartId - it might be populated (object) or just an ObjectId
    let cartId = null;
    if (orderData.cartId) {
      if (typeof orderData.cartId === 'object' && orderData.cartId.cartId) {
        cartId = orderData.cartId.cartId;
      } else if (mongoose.Types.ObjectId.isValid(orderData.cartId)) {
        cartId = orderData.cartId.toString();
      } else {
        cartId = orderData.cartId;
      }
    }
    
    return new OrderDto({
      orderId: orderData.orderId || orderData._id?.toString(),
      orderDate: orderData.orderDate,
      orderDesc: orderData.orderDesc,
      orderFee: orderData.orderFee ? parseFloat(orderData.orderFee) : null,
      cart: cartId ? new CartDto({
        cartId: cartId
      }) : null
    });
  }

  static mapToEntity(orderDto) {
    if (!orderDto) return null;
    
    const orderData = {
      orderDate: orderDto.orderDate,
      orderDesc: orderDto.orderDesc,
      orderFee: orderDto.orderFee
    };
    
    // Convert cartId to ObjectId if provided
    if (orderDto.cart?.cartId || orderDto.cartId) {
      const cartId = orderDto.cart?.cartId || orderDto.cartId;
      if (mongoose.Types.ObjectId.isValid(cartId)) {
        orderData.cartId = new mongoose.Types.ObjectId(cartId);
      } else {
        orderData.cartId = cartId;
      }
    }
    
    // Remove undefined values
    Object.keys(orderData).forEach(key => {
      if (orderData[key] === undefined) {
        delete orderData[key];
      }
    });
    
    return orderData;
  }
}

export default OrderMappingHelper;
