import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";
import { getCartByUserId, clearCartByUserId } from "../services/cartService.js"; // ✅ use service functions

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Load proto
const packageDef = protoLoader.loadSync(
  path.join(__dirname, "./proto/cart.proto")
);
const cartProto = grpc.loadPackageDefinition(packageDef).cart;

const server = new grpc.Server();

// Map gRPC methods to service functions
server.addService(cartProto.CartService.service, {
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

const PORT = 5001;
server.bindAsync(
  `0.0.0.0:${PORT}`,
  grpc.ServerCredentials.createInsecure(),
  () => {
    console.log(`Cart gRPC server running on port ${PORT}`);
    server.start();
  }
);
