// src/controllers/order.controller.js
import Order from "../models/order.model.js";
import { getCartGrpc, clearCartGrpc } from "../grpc/clients/cart.client.js";
import eventService from "../services/event.service.js";
import { createPaymentGrpc } from "../grpc/clients/paymentClient.js";

// -------------------------------
// Create a new order from user's cart
// -------------------------------
export const createOrder = async (req, res, next) => {
  try {
    const userId = 123

    const cart = await getCartGrpc(userId);
    if (!cart || cart.items.length === 0) {
      return res.status(400).json({ message: "Cart is empty" });
    }

    const totalAmount = cart.items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );
// todo  for test purpose only
const orderNumber = "ORD-" + Date.now(); // or any unique logic

    const order = await Order.create({
      userId,
      items: cart.items,
      status: "CREATED",
      orderNumber
    });

    // 🔹 gRPC call to Payment Service
    const payment = await createPaymentGrpc(
      order._id.toString(),
      userId,
      totalAmount
      
    );

    if (payment.status !== "SUCCESS") {
      order.status = "PAYMENT_FAILED";
      await order.save();
      return res.status(402).json({ message: "Payment failed" });
    }

    order.status = "PAID";
    await order.save();

    await clearCartGrpc(userId);

    res.status(201).json(order);
  } catch (error) {
    next(error);
  }
};

// -------------------------------
// Get order by ID
// -------------------------------
export const getOrderById = async (req, res, next) => {
  try {
    const order = await Order.findById(req.params.id);
    if (!order) {
      return res.status(404).json({ message: "Order not found" });
    }
    res.status(200).json(order);
  } catch (error) {
    next(error);
  }
};

// -------------------------------
// Get all orders for a user
// -------------------------------
export const getOrdersByUser = async (req, res, next) => {
    // todo: to be remove 
  console.log("Update order status called");
  try {
    // const userId = req.user.id;
    // TOdo : temporary userId for testing
    const userId = req.params.userId;
    const orders = await Order.find({ userId }).sort({ createdAt: -1 });
    res.status(200).json(orders);
  } catch (error) {
    next(error);
  }
};

// -------------------------------
// Update order status manually
// -------------------------------
export const updateOrderStatus = async (req, res, next) => {

  try {
    const { status } = req.body;
    const order = await Order.findByIdAndUpdate(
      req.params.id,
      { status },
      { new: true }
    );
    if (!order) {
      return res.status(404).json({ message: "Order not found" });
    }

    await eventService.publish("ORDER_UPDATED", order);

    res.status(200).json(order);
  } catch (error) {
    next(error);
  }
};

// -------------------------------
// Cancel/Delete an order
// -------------------------------
export const cancelOrder = async (req, res, next) => {
  try {
    const order = await Order.findByIdAndDelete(req.params.id);
    if (!order) {
      return res.status(404).json({ message: "Order not found" });
    }

    await eventService.publish("ORDER_CANCELLED", order);

    res.status(204).send();
  } catch (error) {
    next(error);
  }
};
