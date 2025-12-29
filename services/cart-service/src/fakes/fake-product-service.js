import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const packageDef = protoLoader.loadSync(
  path.join(__dirname, "../proto/product.proto")
);

const productProto = grpc.loadPackageDefinition(packageDef).product;

// FAKE DATA
const fakeProducts = {
  "p1": { id: "p1", name: "Laptop", price: 1200, stock: 3 },
  "p2": { id: "p2", name: "Mouse", price: 25, stock: 100 },
};

// Fake implementation
function GetProduct(call, callback) {
  const product = fakeProducts[call.request.productId];

  if (!product) {
    return callback({
      code: grpc.status.NOT_FOUND,
      message: "Product not found",
    });
  }

  callback(null, product);
}

// Start fake server
const server = new grpc.Server();
server.addService(productProto.ProductService.service, { GetProduct });

server.bindAsync(
  "0.0.0.0:5002",
  grpc.ServerCredentials.createInsecure(),
  () => {
    console.log("🧪 Fake Product Service running on port 5002");
    server.start();
  }
);
