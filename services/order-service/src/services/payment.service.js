import axios from "axios";
import dotenv from "dotenv";
dotenv.config();
// Base URL of the Payment Service
const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || "http://localhost:5003";
// Create a payment for an order
const createPayment = async (orderId, amount) => {
  try {
    const response = await axios.post(`${PAYMENT_SERVICE_URL}/payments`, {
      orderId,
      amount,
    });

    return response.data; // Expected to return { status: 'PAID' } or { status: 'FAILED' }
  } catch (error) {
    console.error(
      `Payment creation failed for order ${orderId}`,
      error.response?.data || error.message
    );
    return { status: "FAILED" };
  }
};
const getPaymentStatus = async (paymentId) => {
  try {
    const response = await axios.get(`${PAYMENT_SERVICE_URL}/payments/${paymentId}`);
    return response.data; // Expected { status: 'PAID' | 'FAILED' | 'PENDING' }
  } catch (error) {
    console.error(`Failed to fetch payment status ${paymentId}`, error.response?.data || error.message);
    return { status: "UNKNOWN" };
  }
};

export default {
  createPayment,
  getPaymentStatus,
};
