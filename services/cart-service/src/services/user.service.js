import { productClient } from "../grpc/grpc-client.js";

export const getUser = (userId) => {
  return new Promise((resolve, reject) => {
    userClient.GetUser({ userId }, (err, response) => {
      if (err) return reject(err);
      resolve(response);
    });
  });
};
