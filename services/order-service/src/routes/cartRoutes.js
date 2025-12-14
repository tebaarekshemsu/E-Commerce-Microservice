import express from "express";
const router = express.Router();
import { body, param } from "express-validator";
import cartController from "../controllers/CartController.js";

// Validation rules
const cartIdValidation = param("cartId")
  .notEmpty()
  .withMessage("Input must not be blank")
  .isInt()
  .withMessage("Cart ID must be an integer");

const cartBodyValidation = [
  body("userId").optional().isInt().withMessage("User ID must be an integer"),
];

// Routes
router.get("/", cartController.findAll.bind(cartController));

router.get(
  "/:cartId",
  cartIdValidation,
  cartController.findById.bind(cartController)
);

router.post(
  "/",
  body().notEmpty().withMessage("Input must not be NULL!"),
  ...cartBodyValidation,
  cartController.save.bind(cartController)
);

router.put(
  "/",
  body().notEmpty().withMessage("Input must not be NULL"),
  ...cartBodyValidation,
  cartController.update.bind(cartController)
);

router.put(
  "/:cartId",
  cartIdValidation,
  body().notEmpty().withMessage("Input must not be NULL"),
  ...cartBodyValidation,
  cartController.updateById.bind(cartController)
);

router.delete("/:cartId", cartController.deleteById.bind(cartController));

export default router;
