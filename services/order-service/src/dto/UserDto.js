class UserDto {
  constructor(data = {}) {
    this.userId = data.userId || null;
    this.firstName = data.firstName || null;
    this.lastName = data.lastName || null;
    this.imageUrl = data.imageUrl || null;
    this.email = data.email || null;
    this.phone = data.phone || null;
    this.cart = data.cart || null;
  }

  toJSON() {
    const json = {
      userId: this.userId,
      firstName: this.firstName,
      lastName: this.lastName,
      imageUrl: this.imageUrl,
      email: this.email,
      phone: this.phone
    };
    
    if (this.cart) {
      json.cart = this.cart;
    }
    
    return json;
  }
}

export default UserDto;

