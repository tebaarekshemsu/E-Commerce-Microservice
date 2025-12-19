// src/controllers/cart.controller.js
import Cart from "../models/cart.model.js";
import productService from "../services/product.service.js";

// Add item to cart
export const addItem = async (req, res, next) => {
  try {
    const { productId, quantity } = req.body;
    const userId = req.user.id;

    // Validate quantity
    if (!quantity || quantity <= 0) {
      return res.status(400).json({ message: "Quantity must be greater than 0" });
    }

    // Get product details from Product Service
    const product = await productService.getProduct(productId);
    if (!product) {
      return res.status(404).json({ message: "Product not found" });
    }

    // Check stock availability
    if (product.stock < quantity) {
      return res.status(400).json({ message: "Insufficient stock" });
    }

    // Add or update item in cart
    const cart = await Cart.findOneAndUpdate(
      { userId, "items.productId": { $ne: productId } },
      {
        $push: {
          items: {
            productId,
            quantity,
            price: product.price,
          },
        },
      },
      { new: true, upsert: true }
    );

    res.status(200).json(cart);
  } catch (error) {
    next(error);
  }
};

// Get cart by userId
export const getCart = async (req, res, next) => {
  try {
    const userId = req.params.userId || req.user.id;

    const cart = await Cart.findOne({ userId });
    res.status(200).json(cart || { userId, items: [] });
  } catch (error) {
    next(error);
  }
};

// Update item quantity
export const updateItemQuantity = async (req, res, next) => {
  try {
    const { productId, quantity } = req.body;
    const userId = req.user.id;

    // Validate quantity
    if (!quantity || quantity <= 0) {
      return res.status(400).json({ message: "Quantity must be greater than 0" });
    }

    // Validate product existence and stock
    const product = await productService.getProduct(productId);
    if (!product) {
      return res.status(404).json({ message: "Product not found" });
    }
    if (product.stock < quantity) {
      return res.status(400).json({ message: "Insufficient stock" });
    }

    // Update quantity in cart
    const cart = await Cart.findOneAndUpdate(
      { userId, "items.productId": productId },
      { $set: { "items.$.quantity": quantity } },
      { new: true }
    );

    if (!cart) {
      return res.status(404).json({ message: "Cart or item not found" });
    }

    res.status(200).json(cart);
  } catch (error) {
    next(error);
  }
};

// Remove item from cart
export const removeItem = async (req, res, next) => {
  try {
    const { productId } = req.params;
    const userId = req.user.id;

    const cart = await Cart.findOneAndUpdate(
      { userId },
      { $pull: { items: { productId } } },
      { new: true }
    );

    if (!cart) {
      return res.status(404).json({ message: "Cart or item not found" });
    }

    res.status(200).json(cart);
  } catch (error) {
    next(error);
  }
};

// Clear entire cart
export const clearCart = async (req, res, next) => {
  try {
    const userId = req.user.id;

    await Cart.findOneAndDelete({ userId });
    res.status(204).send();
  } catch (error) {
    next(error);
  }
};
