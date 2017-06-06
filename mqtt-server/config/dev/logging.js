const winston = require('winston');

module.exports = {
  transports: [
    new winston.transports.Console({
      level: 'debug',
      colorize: true,
      timestamp: true,
      json: true,
      stringify: true
    })
  ]
};