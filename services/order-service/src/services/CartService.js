import axios from 'axios';
import { Cart, Order } from '../models/index.js';
import CartDto from '../dto/CartDto.js';
import UserDto from '../dto/UserDto.js';
import CartMappingHelper from '../helpers/CartMappingHelper.js';
import CartNotFoundException from '../exceptions/CartNotFoundException.js';
import constants from '../config/constants.js';
import mongoose from 'mongoose';

const { DISCOVERED_DOMAINS_API } = constants;

class CartService {
  async findAll() {
    console.log('*** CartDto List, service; fetch all carts *');
    const carts = await Cart.find().exec();
    const cartDtos = carts.map(cart => CartMappingHelper.mapToDto(cart));
    
    // Fetch user details for each cart
    const cartsWithUsers = await Promise.all(
      cartDtos.map(async (cartDto) => {
        if (cartDto.userId) {
          try {
            const userResponse = await axios.get(
              `${DISCOVERED_DOMAINS_API.USER_SERVICE_API_URL}/${cartDto.userId}`
            );
            cartDto.user = new UserDto(userResponse.data);
          } catch (error) {
            console.error(`Error fetching user ${cartDto.userId}:`, error.message);
            // Continue without user data if fetch fails
          }
        }
        return cartDto;
      })
    );
    
    return cartsWithUsers;
  }

  async findById(cartId) {
    console.log('*** CartDto, service; fetch cart by id *');
    if (!mongoose.Types.ObjectId.isValid(cartId)) {
      throw new CartNotFoundException(`Cart with id: ${cartId} not found`);
    }
    const cart = await Cart.findById(cartId).exec();
    if (!cart) {
      throw new CartNotFoundException(`Cart with id: ${cartId} not found`);
    }
    const cartDto = CartMappingHelper.mapToDto(cart);
    // Fetch user details
    if (cartDto.userId) {
      try {
        const userResponse = await axios.get(
          `${DISCOVERED_DOMAINS_API.USER_SERVICE_API_URL}/${cartDto.userId}`
        );
        cartDto.user = new UserDto(userResponse.data);
      } catch (error) {
        console.error(`Error fetching user ${cartDto.userId}:`, error.message);
        // Continue without user data if fetch fails
      }
    }
    
    return cartDto;
  }

  async save(cartDto) {
    console.log('*** CartDto, service; save cart *');
    const cartData = CartMappingHelper.mapToEntity(cartDto);
    const cart = await Cart.create(cartData);
    
    return CartMappingHelper.mapToDto(cart);
  }

  async update(cartDto) {
    console.log('*** CartDto, service; update cart *');
    if (!cartDto.cartId) {
      throw new CartNotFoundException('Cart ID is required for update');
    }
    
    // First verify the cart exists
    await this.findById(cartDto.cartId);
    
    const cartData = CartMappingHelper.mapToEntity(cartDto);
    delete cartData.cartId; // Remove cartId from update data
    const cart = await Cart.findByIdAndUpdate(
      cartDto.cartId,
      { $set: cartData },
      { new: true, runValidators: true }
    ).exec();
    
    return CartMappingHelper.mapToDto(cart);
  }

  async updateById(cartId, cartDto) {
    console.log('*** CartDto, service; update cart with cartId *');
    // First verify the cart exists
    await this.findById(cartId);
    
    const cartData = CartMappingHelper.mapToEntity(cartDto);
    delete cartData.cartId; // Remove cartId from update data
    
    const cart = await Cart.findByIdAndUpdate(
      cartId,
      { $set: cartData },
      { new: true, runValidators: true }
    ).exec();
    
    return CartMappingHelper.mapToDto(cart);
  }

  async deleteById(cartId) {
    console.log('*** Void, service; delete cart by id *');
    await this.findById(cartId); // Verify cart exists
    
    await Cart.findByIdAndDelete(cartId);
    
    return true;
  }
}

export default new CartService();
