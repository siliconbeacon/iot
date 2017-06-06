const mqtt = require('./mqtt');
const winston = require('winston');
const config = require('./config');

winston.configure(config.logging);
const mqttServer = new mqtt.Server(config.mqtt);
