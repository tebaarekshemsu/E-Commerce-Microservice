// src/controllers/cart.controller.js
import Cart from "../models/cart.model.js";
import * as productService from "../services/product.service.js";
import * as userService from "../services/user.service.js"; // Optional, for user validation

// Add item to cart
export const addItem = async (req, res, next) => {
  try {
    const { productId, quantity } = req.body;
    // I want to know  type of quantity
    console.log(typeof quantity);
    // change to number
    const qty = Number(quantity);
    // const userId = req.user.id;

    // Validate quantity
    if (!qty || qty <= 0) {
      return res
        .status(400)
        .json({ message: "Quantity must be greater than 0" });
    }

    // Validate user via gRPC (optional)
    // const user = await userService.getUser(userId);
    // if (!user || !user.active) return res.status(404).json({ message: "User not found or inactive" });

    // Get product details from Product Service via gRPC
    const product = await productService.getProduct(productId);
    if (!product) {
      return res.status(404).json({ message: "Product not found" });
    }

    // Check stock availability
    if (product.stock < qty) {
      return res.status(400).json({ message: "Insufficient stock" });
    }
    const userId = 123; // Temporary hardcoded userId for testing
    // Add or update item in cart
    let cart = await Cart.findOne({ userId });
    // TOdodo: userID for test

    if (!cart) {
      // Create new cart if it doesn't exist
      cart = new Cart({
        userId: 123, // Temporary hardcoded userId for testing
        items: [
          {
            productId,
            quantity: qty,
            price: product.price,
            name: product.name,
            image: product.image,
          },
        ],
        totalQuantity: qty,
        totalPrice: qty * product.price,
      });
    } else {
      // Update existing cart
      const existingItemIndex = cart.items.findIndex(
        (item) => item.productId === productId
      );
      if (existingItemIndex > -1) {
        // Update quantity
        cart.items[existingItemIndex].quantity += qty;
      } else {
        // Add new item
        cart.items.push({
          productId,
          quantity: qty,
          price: product.price,
          name: product.name,
          image: product.image,
        });
      }

      // Recalculate totals
      cart.totalQuantity = cart.items.reduce(
        (acc, item) => acc + item.quantity,
        0
      );
      cart.totalPrice = cart.items.reduce(
        (acc, item) => acc + item.quantity * item.price,
        0
      );
    }

    await cart.save();

    res.status(200).json(cart);
  } catch (error) {
    next(error);
  }
};

// Get cart by userId
export const getCart = async (req, res, next) => {
  try {
    // const userId = req.params.userId || req.user.id ; // Temporary hardcoded userId for testing
// Todo : to be replaced with userID from auth middleware
    const userId = 123;
    const cart = await Cart.findOne({ userId });
    res
      .status(200)
      .json(cart || { userId, items: [], totalQuantity: 0, totalPrice: 0 });
  } catch (error) {
    next(error);
  }
};

// Update item quantity
export const updateItemQuantity = async (req, res, next) => {
  try {
    const { productId, quantityChange } = req.body;
    const userId = 123; //  todo:Hardcoded for testing
    const qtyChange = Number(quantityChange);
    // if (!quantityChange || isNaN(quantityChange))
    //   return res
    //     .status(400)
    //     .json({ message: "Quantity change must be a number" });

    const cart = await Cart.findOne({ userId });
    if (!cart) return res.status(404).json({ message: "Cart not found" });

    const itemIndex = cart.items.findIndex(
      (item) => item.productId === productId
    );
    if (itemIndex === -1)
      return res.status(404).json({ message: "Item not found in cart" });

    const product = await productService.getProduct(productId);
    if (!product) return res.status(404).json({ message: "Product not found" });

    const newQty = cart.items[itemIndex].quantity + qtyChange;

    if (newQty > product.stock)
      return res.status(400).json({ message: "Insufficient stock" });

    if (newQty <= 0) {
      // Remove item from cart
      cart.items.splice(itemIndex, 1);
    } else {
      cart.items[itemIndex].quantity = newQty;
    }

    // Recalculate totals
    cart.totalQuantity = cart.items.reduce(
      (acc, item) => acc + item.quantity,
      0
    );
    cart.totalPrice = cart.items.reduce(
      (acc, item) => acc + item.quantity * item.price,
      0
    );

    await cart.save();

    res.status(200).json(cart);
  } catch (error) {
    next(error);
  }
};

// Remove item from cart
export const removeItem = async (req, res, next) => {
  try {
    const { productId } = req.params;
    // type of produtId
    console.log(typeof productId);
    // const userId = req.user.id;
    // Todo : to be replaced with userID from auth middleware
    const userId = 123; // Temporary hardcoded userId for testing
    const cart = await Cart.findOne({ userId });
    if (!cart) return res.status(404).json({ message: "Cart not found" });
    cart.items = cart.items.filter((item) => item.productId !== productId);
    // todo to be removed after test
    console.log(cart.items);
    // Recalculate totals
    cart.totalQuantity = cart.items.reduce(
      (acc, item) => acc + item.quantity,
      0
    );
    cart.totalPrice = cart.items.reduce(
      (acc, item) => acc + item.quantity * item.price,
      0
    );

    await cart.save();

    res.status(200).json(cart);
  } catch (error) {
    next(error);
  }
};

// Clear entire cart
export const clearCart = async (req, res, next) => {
  try {
    // const userId = req.user.id;
    // Todo : to be replaced with userID from auth middleware
    const userId = 123; // Temporary hardcoded userId for testing

    await Cart.findOneAndDelete({ userId });
    res.status(200).send({ message: "Cart cleared successfully" });
  } catch (error) {
    next(error);
  }
};

// get all carts 
export const getAllCarts = async (req, res, next) => {
  try {
    const carts = await Cart.find();
    res.status(200).json(carts);
  } catch (error) {
    next(error);
  }
};

