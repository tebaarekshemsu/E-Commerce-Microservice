import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";

const packageDef = protoLoader.loadSync("../grpc/protos/payment.proto");
const proto = grpc.loadPackageDefinition(packageDef).payment;

const createPayment = (call, callback) => {
  const { orderId, userId, amount } = call.request;

  console.log("💳 Fake Payment Request:", {
    orderId,
    userId,
    amount,
  });

  // Fake logic
  const shouldFail = false; // toggle for testing

  if (shouldFail) {
    return callback(null, {
      status: "FAILED",
      paymentId: "",
    });
  }

  callback(null, {
    status: "SUCCESS",
    paymentId: "fake_pay_" + Date.now(),
  });
};

const server = new grpc.Server();
server.addService(proto.PaymentService.service, {
  CreatePayment: createPayment,
});

server.bindAsync(
  "0.0.0.0:5004",
  grpc.ServerCredentials.createInsecure(),
  () => {
    console.log("✅ Fake Payment gRPC running on port 5004");
    server.start();
  }
);
