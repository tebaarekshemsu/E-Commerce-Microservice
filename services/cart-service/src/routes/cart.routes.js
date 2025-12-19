// src/routes/cart.routes.js
import express from "express";
import {
  addItem,
  getCart,
  updateItemQuantity,
  removeItem,
  clearCart,
} from "../controllers/cart.controller.js";
import authMiddleware from "../middlewares/auth.middleware.js";

const router = express.Router();

// Add item to cart
router.post("/", authMiddleware, addItem);

// Get cart by userId (or logged-in user)
router.get("/:userId?", authMiddleware, getCart);

// Update quantity of an item in the cart
router.put("/", authMiddleware, updateItemQuantity);

// Remove a single item from the cart
router.delete("/item/:productId", authMiddleware, removeItem);

// Clear entire cart for the user
router.delete("/clear", authMiddleware, clearCart);

export default router;
