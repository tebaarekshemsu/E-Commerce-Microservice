import OrderNotFoundException from '../exceptions/OrderNotFoundException.js';
import CartNotFoundException from '../exceptions/CartNotFoundException.js';

const errorHandler = (err, req, res, next) => {
  console.error('Error:', err);
  
  // Handle validation errors (Mongoose validation errors)
  if (err.name === 'ValidationError' || err.name === 'CastError') {
    return res.status(400).json({
      msg: `*${err.message}!**`,
      httpStatus: 400,
      timestamp: new Date().toISOString()
    });
  }
  
  // Handle custom exceptions
  if (err instanceof OrderNotFoundException || 
      err instanceof CartNotFoundException ||
      err.name === 'OrderNotFoundException' ||
      err.name === 'CartNotFoundException' ||
      err instanceof Error && err.message.includes('not found')) {
    return res.status(400).json({
      msg: `#### ${err.message}! ####`,
      httpStatus: 400,
      timestamp: new Date().toISOString()
    });
  }
  
  // Handle syntax errors (JSON parsing)
  if (err instanceof SyntaxError && err.status === 400 && 'body' in err) {
    return res.status(400).json({
      msg: '*Invalid JSON format!**',
      httpStatus: 400,
      timestamp: new Date().toISOString()
    });
  }
  
  // Handle illegal state errors
  if (err instanceof Error && err.message.includes('IllegalState')) {
    return res.status(400).json({
      msg: `#### ${err.message}! ####`,
      httpStatus: 400,
      timestamp: new Date().toISOString()
    });
  }
  
  // Default error handler
  res.status(err.statusCode || 500).json({
    msg: err.message || 'Internal server error',
    httpStatus: err.statusCode || 500,
    timestamp: new Date().toISOString()
  });
};

export default errorHandler;

