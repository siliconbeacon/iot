const mosca = require('mosca');
const redis = require('redis');
const winston = require('winston');
const config = require('./config');

const mqttClientSet = 'mqtt:clients';

function Server(options) {
  const _mosca = this._mosca = new mosca.Server(options);
  const _redisClient = this._redisClient = redis.createClient(config.redis);
  _mosca.on('ready', () => {
    _redisClient.del(mqttClientSet);
    winston.info("MQTT/S Listening on port " + _mosca.opts.secure.port);
    winston.info("MQTT over WebSockets listening on port " + _mosca.opts.https.port);
  });
  _mosca.on('clientConnected', client => {
    winston.info('Client connected: ' + client.id);
    _redisClient.sadd(mqttClientSet, client.id);
  });
  _mosca.on('clientDisconnected', client => {
    winston.info('Client disconnected: ' + client.id);
    _redisClient.srem(mqttClientSet, client.id);
  });
}

exports.Server = Server;