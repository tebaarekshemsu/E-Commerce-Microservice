class OrderDto {
  constructor(data = {}) {
    this.orderId = data.orderId || null;
    this.orderDate = data.orderDate || null;
    this.orderDesc = data.orderDesc || null;
    this.orderFee = data.orderFee || null;
    this.cart = data.cart || null;
  }

  toJSON() {
    return {
      orderId: this.orderId,
      orderDate: this.orderDate,
      orderDesc: this.orderDesc,
      orderFee: this.orderFee,
      ...(this.cart && { cart: this.cart })
    };
  }
}

export default OrderDto;

