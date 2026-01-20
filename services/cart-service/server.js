import dotenv from "dotenv";
dotenv.config();

import express from "express";
import cors from "cors";
import morgan from "morgan";
import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

import cartroutes from "./src/routes/cart.routes.js";
import { connectDB } from "./src/config/database.js";
import { getCartByUserId, clearCartByUserId } from "./src/services/cartService.js";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
const REST_PORT = process.env.PORT || 5001;
const GRPC_PORT = oprocess.env.CART_GRPC_PORT||5003; // Separate port for gRPC

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(morgan("dev"));

// Health check
app.get("/", (req, res) => {
  res.json({ message: "Cart REST Service running!" });
});

// REST routes
app.use("/cart_service/api/cart", cartroutes);

// Start REST + gRPC after DB connected
const startServer = async () => {
  try {
    await connectDB();

    // Start REST API
    app.listen(REST_PORT, () => {
      console.log(`Cart REST Service running on port ${REST_PORT}`);
    });

    // gRPC setup
    const packageDef = protoLoader.loadSync(path.join(__dirname, "./src/grpc/proto/cart.proto"));
    const cartProto = grpc.loadPackageDefinition(packageDef).cart;

    const grpcServer = new grpc.Server();

    grpcServer.addService(cartProto.CartService.service, {
      GetCart: async (call, callback) => {
        try {
          const cart = await getCartByUserId(call.request.userId);
          callback(null, cart);
        } catch (err) {
          callback(err, null);
        }
      },
      ClearCart: async (call, callback) => {
        try {
          const result = await clearCartByUserId(call.request.userId);
          callback(null, result);
        } catch (err) {
          callback(err, null);
        }
      },
    });

    grpcServer.bindAsync(
      `0.0.0.0:${GRPC_PORT}`,
      grpc.ServerCredentials.createInsecure(),
      (err, port) => {
        if (err) {
          console.error("Failed to bind gRPC server:", err);
          return;
        }
        grpcServer.start();
        console.log(`Cart gRPC server running on port ${port}`);
      }
    );

  } catch (error) {
    console.error("Unable to start Cart Service:", error);
    process.exit(1);
  }
};

startServer();

export default app;
