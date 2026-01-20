// src/grpc-client.js
import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";
// Fix __dirname for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
// Load Product Proto
const productPackageDef = protoLoader.loadSync(
  path.join(__dirname, "./proto/product.proto")
);
const productProto = grpc.loadPackageDefinition(productPackageDef).product;
export const productClient = new productProto.ProductService(
  process.env.PRODUCT_SERVICE_URL || "localhost:5004",
  grpc.credentials.createInsecure()
);

// Load User Proto
const userPackageDef = protoLoader.loadSync(
  path.join(__dirname, "./proto/user.proto")
);
const userProto = grpc.loadPackageDefinition(userPackageDef).user;

export const userClient = new userProto.UserService(
  process.env.USER_SERVICE_URL || "localhost:5002",
  grpc.credentials.createInsecure()
);
