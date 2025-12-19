// cart.model.js
import mongoose from "mongoose";

const cartItemSchema = new mongoose.Schema(
  {
	productId: { type: String, required: true },
	quantity: { type: Number, required: true },

	// snapshot (important!)
	price: { type: Number, required: true },
	name: { type: String },
	image: { type: String },
  },
  { _id: false }
);

const cartSchema = new mongoose.Schema(
  {
	userId: {
	  type: String,
	  required: true,
	  index: true,
	  unique: true,
	},

	items: [cartItemSchema],

	totalQuantity: { type: Number, default: 0 },
	totalPrice: { type: Number, default: 0 },

	expiresAt: {
	  type: Date,
	  index: { expireAfterSeconds: 0 },
	},
  },
  { timestamps: true }
);

export default mongoose.model("Cart", cartSchema);
