import { productClient } from "../grpc/grpc-client.js";

export const getProduct = (productId) => {
  return new Promise((resolve, reject) => {
    productClient.GetProduct({ productId }, (err, response) => {
      if (err) return reject(err);
      resolve(response);
    });
  });
};
