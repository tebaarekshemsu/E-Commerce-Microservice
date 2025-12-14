import CartDto from '../dto/CartDto.js';
import UserDto from '../dto/UserDto.js';

class CartMappingHelper {
  static mapToDto(cart) {
    if (!cart) return null;
    
    // Convert Mongoose document to plain object
    const cartData = cart.toObject ? cart.toObject() : cart;
    
    return new CartDto({
      cartId: cartData.cartId || cartData._id?.toString(),
      userId: cartData.userId,
      user: cartData.userId ? new UserDto({
        userId: cartData.userId
      }) : null
    });
  }

  static mapToEntity(cartDto) {
    if (!cartDto) return null;
    
    const cartData = {
      userId: cartDto.userId
    };
    
    // Remove undefined values
    Object.keys(cartData).forEach(key => {
      if (cartData[key] === undefined) {
        delete cartData[key];
      }
    });
    
    return cartData;
  }
}

export default CartMappingHelper;
