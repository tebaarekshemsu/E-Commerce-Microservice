import express from 'express';
const router = express.Router();
import { body, param } from 'express-validator';
import orderController from '../controllers/OrderController.js';

// Validation rules
const orderIdValidation = param('orderId')
  .notEmpty()
  .withMessage('Input must not be blank')
  .isInt()
  .withMessage('Order ID must be an integer');

const orderBodyValidation = [
  body('orderDate').optional().isISO8601().withMessage('Invalid date format'),
  body('orderDesc').optional().isString().withMessage('Order description must be a string'),
  body('orderFee').optional().isNumeric().withMessage('Order fee must be a number'),
  body('cart').optional().isObject().withMessage('Cart must be an object'),
  body('cart.cartId').optional().isInt().withMessage('Cart ID must be an integer')
];

// Routes
router.get('/', orderController.findAll.bind(orderController));

router.get('/:orderId', 
  orderIdValidation,
  orderController.findById.bind(orderController)
);

router.post('/', 
  body().notEmpty().withMessage('Input must not be NULL'),
  ...orderBodyValidation,
  orderController.save.bind(orderController)
);

router.put('/', 
  body().notEmpty().withMessage('Input must not be NULL'),
  ...orderBodyValidation,
  orderController.update.bind(orderController)
);

router.put('/:orderId', 
  orderIdValidation,
  body().notEmpty().withMessage('Input must not be NULL'),
  ...orderBodyValidation,
  orderController.updateById.bind(orderController)
);

router.delete('/:orderId', 
  orderController.deleteById.bind(orderController)
);

export default router;

