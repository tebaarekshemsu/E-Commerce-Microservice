import cartService from '../services/CartService.js';
import { body, param, validationResult } from 'express-validator';

class CartController {
  async findAll(req, res, next) {
    try {
      console.log('*** CartDto List, controller; fetch all categories *');
      const carts = await cartService.findAll();
      res.json({ collection: carts });
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
      
      console.log('*** CartDto, resource; fetch cart by id *');
      const cartId = parseInt(req.params.cartId);
      const cart = await cartService.findById(cartId);
      res.json(cart.toJSON());
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
      
      console.log('*** CartDto, resource; save cart *');
      const cart = await cartService.save(req.body);
      res.json(cart.toJSON());
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
      
      console.log('*** CartDto, resource; update cart *');
      const cart = await cartService.update(req.body);
      res.json(cart.toJSON());
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
      
      console.log('*** CartDto, resource; update cart with cartId *');
      const cartId = parseInt(req.params.cartId);
      const cart = await cartService.updateById(cartId, req.body);
      res.json(cart.toJSON());
    } catch (error) {
      next(error);
    }
  }

  async deleteById(req, res, next) {
    try {
      console.log('*** Boolean, resource; delete cart by id *');
      const cartId = parseInt(req.params.cartId);
      await cartService.deleteById(cartId);
      res.json(true);
    } catch (error) {
      next(error);
    }
  }
}

export default new CartController();

