class CartNotFoundException extends Error {
  constructor(message) {
    super(message);
    this.name = 'CartNotFoundException';
    this.statusCode = 400;
  }
}

export default CartNotFoundException;

