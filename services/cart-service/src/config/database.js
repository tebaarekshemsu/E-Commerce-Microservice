import dotenv from "dotenv";
dotenv.config();
import mongoose from "mongoose";

const connectDB = async () => {
  console.log("Connecting to MongoDB...", process.env.MONGO_URI);
  try {
    const mongoURI =
      process.env.MONGO_URI ||
      `mongodb://${process.env.DB_HOST || "localhost"}:${
        process.env.DB_PORT || "27017"
      }/${process.env.DB_NAME || "ecommerce_dev_db"}`;

    const options = {
      useNewUrlParser: true,
      useUnifiedTopology: true,
    };

    if (process.env.DB_USER && process.env.DB_PASSWORD) {
      options.auth = {
        username: process.env.DB_USER,
        password: process.env.DB_PASSWORD,
      };
    }

    await mongoose.connect(mongoURI, options);
    console.log("MongoDB connected successfully");
  } catch (error) {
    console.error("MongoDB connection error:", error);
    process.exit(1);
  }
};

export { connectDB, mongoose };
