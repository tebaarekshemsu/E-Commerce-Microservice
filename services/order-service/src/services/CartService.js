// src/services/cart.service.js
import axios from "axios";
import dotenv from "dotenv";

dotenv.config();

const CART_SERVICE_URL = process.env.CART_SERVICE_URL || "http://localhost:5002";

// Get user cart from Cart Service
const getCart = async (userId) => {
  try {
    const response = await axios.get(`${CART_SERVICE_URL}/cart/${userId}`);
    return response.data;
  } catch (error) {
    console.error("Failed to fetch cart", error.response?.data || error.message);
    throw new Error("Could not fetch cart");
  }
};

// Clear user cart in Cart Service
const clearCart = async (userId) => {
  try {
    await axios.delete(`${CART_SERVICE_URL}/cart/clear`, {
      headers: { "X-User-Id": userId } // or use auth token
    });
  } catch (error) {
    console.error("Failed to clear cart", error.response?.data || error.message);
    throw new Error("Could not clear cart");
  }
};

// Optional: Add item to cart via Cart Service
const addItem = async (userId, productId, quantity) => {
  try {
    const response = await axios.post(
      `${CART_SERVICE_URL}/cart`,
      { productId, quantity },
      { headers: { "X-User-Id": userId } }
    );
    return response.data;
  } catch (error) {
    console.error("Failed to add item to cart", error.response?.data || error.message);
    throw new Error("Could not add item to cart");
  }
};

// Optional: Remove item from cart via Cart Service
const removeItem = async (userId, productId) => {
  try {
    const response = await axios.delete(`${CART_SERVICE_URL}/cart/item/${productId}`, {
      headers: { "X-User-Id": userId }
    });
    return response.data;
  } catch (error) {
    console.error("Failed to remove item from cart", error.response?.data || error.message);
    throw new Error("Could not remove item from cart");
  }
};

export default {
  getCart,
  clearCart,
  addItem,
  removeItem,
};
