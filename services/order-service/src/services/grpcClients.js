import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

class GRPCClients {
  constructor() {
    this.productClient = null;
    this.paymentClient = null;
    this.userClient = null;
    this.initializeClients();
  }

  initializeClients() {
    try {
      // Product Service Client
      const productProtoPath = path.join(__dirname, "../protos/product.proto");
      const productPackageDefinition = protoLoader.loadSync(productProtoPath, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
      });
      const productProto = grpc.loadPackageDefinition(
        productPackageDefinition,
      ).product;

      const productServiceUrl =
        process.env.PRODUCT_SERVICE_GRPC_URL || "product-service:9001";
      this.productClient = new productProto.ProductService(
        productServiceUrl,
        grpc.credentials.createInsecure(),
      );

      // Payment Service Client
      const paymentProtoPath = path.join(__dirname, "../protos/payment.proto");
      const paymentPackageDefinition = protoLoader.loadSync(paymentProtoPath, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
      });
      const paymentProto = grpc.loadPackageDefinition(
        paymentPackageDefinition,
      ).payment;

      const paymentServiceUrl =
        process.env.PAYMENT_SERVICE_GRPC_URL || "payment-service:9002";
      this.paymentClient = new paymentProto.PaymentService(
        paymentServiceUrl,
        grpc.credentials.createInsecure(),
      );

      // User Service Client
      const userProtoPath = path.join(__dirname, "../protos/user.proto");
      const userPackageDefinition = protoLoader.loadSync(userProtoPath, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
      });
      const userProto = grpc.loadPackageDefinition(userPackageDefinition).user;

      const userServiceUrl =
        process.env.USER_SERVICE_GRPC_URL || "user-service:9003";
      this.userClient = new userProto.UserService(
        userServiceUrl,
        grpc.credentials.createInsecure(),
      );

      console.log("✅ gRPC clients initialized");
    } catch (error) {
      console.error("❌ Failed to initialize gRPC clients:", error);
    }
  }

  // Product Service Methods
  async getProduct(productId) {
    return new Promise((resolve, reject) => {
      this.productClient.GetProduct(
        { id: productId.toString() },
        (error, response) => {
          if (error) {
            console.error("gRPC GetProduct error:", error);
            reject(error);
          } else {
            console.log(`✅ Retrieved product ${productId} via gRPC`);
            resolve(response);
          }
        },
      );
    });
  }

  async checkProductAvailability(productId, quantity) {
    return new Promise((resolve, reject) => {
      this.productClient.CheckAvailability(
        {
          product_id: productId.toString(),
          quantity: quantity,
        },
        (error, response) => {
          if (error) {
            console.error("gRPC CheckAvailability error:", error);
            reject(error);
          } else {
            console.log(
              `✅ Checked availability for product ${productId}, quantity ${quantity}: ${response.available}`,
            );
            resolve(response);
          }
        },
      );
    });
  }

  // Payment Service Methods
  async processPayment(orderData) {
    return new Promise((resolve, reject) => {
      const paymentRequest = {
        order_id: orderData.orderId.toString(),
        amount: orderData.orderFee || 0,
        currency: orderData.currency || "USD",
        payment_method: orderData.paymentMethod || "credit_card",
      };

      this.paymentClient.ProcessPayment(paymentRequest, (error, response) => {
        if (error) {
          console.error("gRPC ProcessPayment error:", error);
          reject(error);
        } else {
          console.log(
            `✅ Processed payment for order ${orderData.orderId}: ${response.success ? "SUCCESS" : "FAILED"}`,
          );
          resolve(response);
        }
      });
    });
  }

  // User Service Methods
  async getUser(userId) {
    return new Promise((resolve, reject) => {
      this.userClient.GetUser({ id: userId.toString() }, (error, response) => {
        if (error) {
          console.error("gRPC GetUser error:", error);
          reject(error);
        } else {
          console.log(`✅ Retrieved user ${userId} via gRPC`);
          resolve(response);
        }
      });
    });
  }

  // Health check methods
  async healthCheck() {
    const checks = {
      product: false,
      payment: false,
      user: false,
    };

    try {
      // Test product service
      await this.getProduct("1");
      checks.product = true;
    } catch (error) {
      console.warn("Product service health check failed:", error.message);
    }

    try {
      // Test user service
      await this.getUser("1");
      checks.user = true;
    } catch (error) {
      console.warn("User service health check failed:", error.message);
    }

    return checks;
  }
}

// Create singleton instance
const grpcClients = new GRPCClients();

export default grpcClients;
