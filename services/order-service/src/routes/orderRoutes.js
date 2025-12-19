// src/routes/order.routes.js
import express from "express";
import {
  createOrder,
  getOrderById,
  getOrdersByUser,
  updateOrderStatus,
  cancelOrder,
} from "../controllers/OrderController.js";
import authMiddleware from "../middlewares/auth.middleware.js";

const router = express.Router();
// Create a new order
router.post("/", authMiddleware, createOrder);
// Get order by ID
router.get("/:id", authMiddleware, getOrderById);

// Get all orders for the logged-in user
router.get("/", authMiddleware, getOrdersByUser);

// Update order status (admin or system)
router.put("/:id/status", authMiddleware, updateOrderStatus);

// Cancel/Delete an order
router.delete("/:id", authMiddleware, cancelOrder);

export default router;
