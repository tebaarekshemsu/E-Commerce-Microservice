import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

// Resolve __dirname for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// ✅ Correct proto path relative to this file
const PROTO_PATH = path.join(__dirname, "../protos/payment.proto");

// Load proto
const packageDef = protoLoader.loadSync(PROTO_PATH);
const proto = grpc.loadPackageDefinition(packageDef).payment;

// Create gRPC client
const client = new proto.PaymentService(
  "localhost:5004", // adjust port if needed
  grpc.credentials.createInsecure()
);

// Export helper function
export const createPaymentGrpc = (orderId, userId, amount) => {
  return new Promise((resolve, reject) => {
    client.CreatePayment(
      { orderId, userId: String(userId), amount },
      (err, response) => {
        if (err) return reject(err);
        resolve(response);
      }
    );
  });
};
