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
// router.post("/", authMiddleware, createOrder);
router.post("/create", createOrder);
// Get order by ID
// router.get("/:id", authMiddleware, getOrderById);
// router.get("/:id", getOrderById);

// Get all orders for the logged-in user
// router.get("/", authMiddleware, getOrdersByUser);
router.get("/:userId", getOrdersByUser);

// Update order status (admin or system)
// router.put("/:id/status", authMiddleware, updateOrderStatus);
router.put("/:id/status", updateOrderStatus);

// Cancel/Delete an order
// router.delete("/:id", authMiddleware, cancelOrder);
router.delete("/:id", cancelOrder);

export default router;
