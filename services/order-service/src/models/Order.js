import mongoose from 'mongoose';

const orderSchema = new mongoose.Schema({
  orderDate: {
    type: Date,
    default: Date.now
  },
  orderDesc: {
    type: String,
    maxlength: 255
  },
  orderFee: {
    type: Number,
    min: 0
  },
  cartId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Cart',
    required: false
  }
}, {
  timestamps: true,
  collection: 'orders'
});

// Transform _id to orderId for consistency with API
orderSchema.set('toJSON', {
  transform: function(doc, ret) {
    ret.orderId = ret._id.toString();
    delete ret._id;
    delete ret.__v;
    if (ret.cartId && ret.cartId._id) {
      ret.cartId = ret.cartId._id.toString();
    }
    return ret;
  }
});

const Order = mongoose.model('Order', orderSchema);

export default Order;
