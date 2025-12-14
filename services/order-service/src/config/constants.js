export default {
  DATE_FORMATS: {
    LOCAL_DATE: 'dd-MM-yyyy',
    LOCAL_DATE_TIME: 'dd-MM-yyyy__HH:mm:ss:SSSSSS',
    ZONED_DATE_TIME: 'dd-MM-yyyy__HH:mm:ss:SSSSSS',
    INSTANT: 'dd-MM-yyyy__HH:mm:ss:SSSSSS'
  },
  DISCOVERED_DOMAINS_API: {
    USER_SERVICE_HOST: process.env.USER_SERVICE_HOST || 'http://USER-SERVICE/user-service',
    USER_SERVICE_API_URL: process.env.USER_SERVICE_URL || 'http://USER-SERVICE/user-service/api/users',
    PRODUCT_SERVICE_HOST: process.env.PRODUCT_SERVICE_HOST || 'http://PRODUCT-SERVICE/product-service',
    PRODUCT_SERVICE_API_URL: process.env.PRODUCT_SERVICE_API_URL || 'http://PRODUCT-SERVICE/product-service/api/products',
    ORDER_SERVICE_HOST: process.env.ORDER_SERVICE_HOST || 'http://ORDER-SERVICE/order-service',
    ORDER_SERVICE_API_URL: process.env.ORDER_SERVICE_API_URL || 'http://ORDER-SERVICE/order-service/api/orders',
    FAVOURITE_SERVICE_HOST: process.env.FAVOURITE_SERVICE_HOST || 'http://FAVOURITE-SERVICE/favourite-service',
    FAVOURITE_SERVICE_API_URL: process.env.FAVOURITE_SERVICE_API_URL || 'http://FAVOURITE-SERVICE/favourite-service/api/favourites',
    PAYMENT_SERVICE_HOST: process.env.PAYMENT_SERVICE_HOST || 'http://PAYMENT-SERVICE/payment-service',
    PAYMENT_SERVICE_API_URL: process.env.PAYMENT_SERVICE_API_URL || 'http://PAYMENT-SERVICE/payment-service/api/payments',
    SHIPPING_SERVICE_HOST: process.env.SHIPPING_SERVICE_HOST || 'http://SHIPPING-SERVICE/shipping-service',
    SHIPPING_SERVICE_API_URL: process.env.SHIPPING_SERVICE_API_URL || 'http://SHIPPING-SERVICE/shipping-service/api/shippings'
  }
};

