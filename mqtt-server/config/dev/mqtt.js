const redis = require('redis');
const redisConfig = require('./redis');

const backend = {
  type: 'redis',
  redis: redis,
  port: redisConfig.port,
  host: redisConfig.host,
  return_buffers: true
};

module.exports = {
  backend: backend,
  port: 1883,
  https: {
    port: 8443,
    bundle: true
  },
  secure: {
    port: 8883,
    keyPath: __dirname + '/iot.siliconbeacon.com.key.pem',
    certPath: __dirname + '/iot.siliconbeacon.com.cert.pem'
  }
};