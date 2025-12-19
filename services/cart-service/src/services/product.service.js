
import axios from "axios";
import dotenv from "dotenv";

dotenv.config();

// Base URL of the Product Service
const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL || "http://localhost:5001";
// Get product details by ID
const getProduct = async (productId) => {
  try {
    const response = await axios.get(`${PRODUCT_SERVICE_URL}/products/${productId}`);
    return response.data;
  } catch (error) {
    console.error(
      `Failed to fetch product ${productId}`,
      error.response?.data || error.message
    );
    return null; 
  }
};

// Optional: Check stock availability
const checkStock = async (productId, quantity) => {
  const product = await getProduct(productId);
  if (!product) return false;
  return product.stock >= quantity;
};

export default {
  getProduct,
  checkStock,
};
