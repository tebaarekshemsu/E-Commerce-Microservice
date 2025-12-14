import mongoose from 'mongoose';
const cartSchema = new mongoose.Schema({
  userId: {
    type: Number,
    required: false
  }
}, {
  timestamps: true,
  collection: 'carts'
});

// Transform _id to cartId for consistency with API
cartSchema.set('toJSON', {
  transform: function(doc, ret) {
    ret.cartId = ret._id.toString();
    delete ret._id;
    delete ret.__v;
    return ret;
  }
});

const Cart = mongoose.model('Cart', cartSchema);

export default Cart;
