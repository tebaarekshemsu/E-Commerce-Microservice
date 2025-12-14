import orderService from '../services/OrderService.js';
import { body, param, validationResult } from 'express-validator';

class OrderController {
  async findAll(req, res, next) {
    try {
      console.log('*** OrderDto List, controller; fetch all orders *');
      const orders = await orderService.findAll();
      res.json({ collection: orders });
    } catch (error) {
      next(error);
    }
  }

  async findById(req, res, next) {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          msg: `*${errors.array()[0].msg}!**`,
          httpStatus: 400,
          timestamp: new Date().toISOString()
        });
      }
      
      console.log('*** OrderDto, resource; fetch order by id *');
      const orderId = parseInt(req.params.orderId);
      const order = await orderService.findById(orderId);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async save(req, res, next) {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          msg: `*${errors.array()[0].msg}!**`,
          httpStatus: 400,
          timestamp: new Date().toISOString()
        });
      }
      
      console.log('*** OrderDto, resource; save order *');
      const order = await orderService.save(req.body);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async update(req, res, next) {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          msg: `*${errors.array()[0].msg}!**`,
          httpStatus: 400,
          timestamp: new Date().toISOString()
        });
      }
      
      console.log('*** OrderDto, resource; update order *');
      const order = await orderService.update(req.body);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async updateById(req, res, next) {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          msg: `*${errors.array()[0].msg}!**`,
          httpStatus: 400,
          timestamp: new Date().toISOString()
        });
      }
      
      console.log('*** OrderDto, resource; update order with orderId *');
      const orderId = parseInt(req.params.orderId);
      const order = await orderService.updateById(orderId, req.body);
      res.json(order.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async deleteById(req, res, next) {
    try {
      console.log('*** Boolean, resource; delete order by id *');
      const orderId = parseInt(req.params.orderId);
      await orderService.deleteById(orderId);
      res.json(true);
    } catch (error) {
      next(error);
    }
  }
}

export default new OrderController();

