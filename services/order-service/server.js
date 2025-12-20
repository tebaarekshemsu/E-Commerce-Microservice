import dotenv from "dotenv";
dotenv.config();
import express from "express";
import cors from "cors";
import morgan from "morgan";
import { connectDB } from "./src/config/database.js";
import orderRoutes from "./src/routes/orderRoutes.js";
import errorHandler from "./src/middlewares/errorHandler.js";
const app = express();
const PORT = process.env.PORT || 8300;
// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(morgan("dev"));
// Health check endpoint
app.get("/", (req, res) => {
  res.json({ message: "Order controller responding!!" });
});
// Routes
app.use(`order-service/api/orders`, orderRoutes);
// Error handling middleware
app.use(errorHandler);
// Database connection and server start
const startServer = async () => {
  try {
    await connectDB();
    app.listen(PORT, () => {
      console.log(`Order Service is running on port ${PORT}`);
    });
  } catch (error) {
    console.error("Unable to start the server:", error);
    process.exit(1);
  }
};

startServer();

export default app;
