import Cart from "../models/cart.model.js";
// Get cart by userId (returns plain object)
export const getCartByUserId = async (userId) => {
  const cart = await Cart.findOne({ userId });
  return cart || { userId, items: [], totalQuantity: 0, totalPrice: 0 };
};
// Clear cart by userId (returns plain object)
export const clearCartByUserId = async (userId) => {
  await Cart.findOneAndDelete({ userId });
  return { message: "Cart cleared successfully" };
};
