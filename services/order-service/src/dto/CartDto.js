class CartDto {
  constructor(data = {}) {
    this.cartId = data.cartId || null;
    this.userId = data.userId || null;
    this.orderDtos = data.orderDtos || null;
    this.user = data.user || null;
  }

  toJSON() {
    const json = {
      cartId: this.cartId,
      userId: this.userId
    };
    
    if (this.orderDtos) {
      json.orderDtos = this.orderDtos;
    }
    
    if (this.user) {
      json.user = this.user;
    }
    
    return json;
  }
}

export default CartDto;

