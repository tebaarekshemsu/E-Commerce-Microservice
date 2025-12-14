class OrderNotFoundException extends Error {
  constructor(message) {
    super(message);
    this.name = 'OrderNotFoundException';
    this.statusCode = 400;
  }
}

export default OrderNotFoundException;

