// src/controllers/order.controller.js
import Order from "../models/order.model.js";
import cartService from "../services/CartService.js";
import eventService from "../services/event.service.js";
import paymentService from "../services/payment.service.js";
// Create a new order from user's cart and initiate payment
export const createOrder = async (req, res, next) => {
  try {
    const userId = req.user.id;
    // Get user's cart
    const cart = await cartService.getCart(userId);
    if (!cart || cart.items.length === 0) {
      return res.status(400).json({ message: "Cart is empty" });
    }

    // Calculate total amount
    const totalAmount = cart.items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );

    // Create order in DB
    const order = await Order.create({
      userId,
      items: cart.items,
      totalAmount,
      status: "PENDING", // initial status
    });

    // Initiate payment through Payment Service
    const paymentResult = await paymentService.createPayment(
      order.id,
      totalAmount
    );
    if (paymentResult.status === "FAILED") {
      order.status = "PAYMENT_FAILED";
      await order.save();
      await eventService.publish("ORDER_PAYMENT_FAILED", order);
      return res.status(402).json({ message: "Payment failed", order });
    }

    // Update order status to paid
    order.status = "PAID";
    await order.save();

    // Publish ORDER_CREATED event
    await eventService.publish("ORDER_CREATED", order);

    // Clear cart after successful order
    await cartService.clearCart(userId);

    res.status(201).json(order);
  } catch (error) {
    next(error);
  }
};

// Get order by ID
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

// Get all orders for a user
export const getOrdersByUser = async (req, res, next) => {
  try {
    const userId = req.user.id;
    const orders = await Order.find({ userId }).sort({ createdAt: -1 });
    res.status(200).json(orders);
  } catch (error) {
    next(error);
  }
};

// Update order status manually
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

// Cancel/Delete an order
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
