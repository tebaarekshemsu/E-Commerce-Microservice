// src/routes/cart.routes.js
import express from "express";
import {
  addItem,
  getCart,
  updateItemQuantity,
  removeItem,
  clearCart,
  getAllCarts,
} from "../controllers/cart.controller.js";
import authMiddleware from "../middlewares/auth.middleware.js";

const router = express.Router();

// Add item to cart
router.post("/", addItem);

// Get cart by userId (or logged-in user)
router.get("/:userId?", getCart);

// Update quantity of an item in the cart
// router.put("/", authMiddleware, updateItemQuantity);
router.put("/update/:productId", updateItemQuantity);

// Remove a single item from the cart
// router.delete("/item/:productId", authMiddleware, removeItem);
router.delete("/item/:productId", removeItem);

// Clear entire cart for the user
// router.delete("/clear/:userId", authMiddleware, clearCart);
router.delete("/clear/:userId", clearCart);
// get all carts
// router.get("/", authMiddleware, getAllCarts);
router.get("/", getAllCarts);

export default router;
