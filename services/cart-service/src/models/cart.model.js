import mongoose from "mongoose";

const cartItemSchema = new mongoose.Schema(
  {
    productId: { type: String, required: true },
    quantity: { type: Number, required: true, min: 1 },

    price: { type: Number, required: true },
    name: String,
    image: String,
  },
  { _id: false }
);

const cartSchema = new mongoose.Schema(
  {
    userId: {
      type: String,
      required: true,
      immutable: true,
      index: true,
      unique: true,
    },

    items: [cartItemSchema],

    totalQuantity: { type: Number, default: 0 },
    totalPrice: { type: Number, default: 0 },

    status: {
      type: String,
      enum: ["ACTIVE", "CHECKED_OUT", "ABANDONED"],
      default: "ACTIVE",
      index: true,
    },

    expiresAt: {
      type: Date,
      index: { expireAfterSeconds: 0 },
    },
  },
  { timestamps: true }
);

export default mongoose.model("Cart", cartSchema);
