const env = process.env.NODE_ENV || 'dev';

module.exports = {
  mqtt: require('./' + env + '/mqtt'),
  logging: require('./' + env + '/logging'),
  redis: require('./' + env + '/redis')
};
