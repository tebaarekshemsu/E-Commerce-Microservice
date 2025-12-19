// order.model.js
import mongoose from "mongoose";

const orderItemSchema = new mongoose.Schema(
  {
    productId: { type: String, required: true },
    // snapshot
    name: String,
    price: Number,
    image: String,

    quantity: Number,
    subtotal: Number,
  },
  { _id: false }
);

const orderSchema = new mongoose.Schema(
  {
    orderNumber: {
      type: String,
      unique: true,
      index: true,
    },

    userId: { type: String, index: true },

    items: [orderItemSchema],

    pricing: {
      subtotal: Number,
      tax: Number,
      discount: Number,
      shippingFee: Number,
      total: Number,
    },

    status: {
      type: String,
      enum: [
        "CREATED",
        "PAYMENT_PENDING",
        "PAID",
        "SHIPPED",
        "DELIVERED",
        "CANCELLED",
        "REFUNDED",
      ],
      index: true,
    },

    payment: {
      method: String,
      transactionId: String,
      paidAt: Date,
      amount: Number,
    },

    shipping: {
      address: Object,
      carrier: String,
      trackingNumber: String,
    },

    version: {
      type: Number,
      default: 1,
    },
  },
  { timestamps: true }
);

export default mongoose.model("Order", orderSchema);
