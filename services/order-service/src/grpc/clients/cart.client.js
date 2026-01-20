import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Load proto
const packageDef = protoLoader.loadSync(
  path.join(__dirname, "../protos/cart.proto")
);
const cartProto = grpc.loadPackageDefinition(packageDef).cart;

// Create raw gRPC client
export const cartClient = new cartProto.CartService(
  process.env.CART_SERVICE_URL || "localhost:5003",
  grpc.credentials.createInsecure()
);

// -------------------------------
// Promise wrappers for async/await
// -------------------------------
export const getCartGrpc = (userId) =>
  new Promise((resolve, reject) => {
    cartClient.GetCart({ userId }, (err, response) => {
      if (err) reject(err);
      else resolve(response);
    });
  });

export const clearCartGrpc = (userId) =>
  new Promise((resolve, reject) => {
    cartClient.ClearCart({ userId }, (err, response) => {
      if (err) reject(err);
      else resolve(response);
    });
  });
