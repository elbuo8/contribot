
function loggerSingleton () {
  if (loggerSingleton.prototype.logger === undefined) {
    loggerSingleton.prototype.logger = console;
    if (process.env.LOGGLY_TOKEN) {
      loggerSingleton.prototype.logger = require('loggly').createClient({
        token: process.env.LOGGLY_TOKEN,
        subdomain: process.env.LOGGLY_SUBDOMAIN,
        auth: {
          username: process.env.LOGGLY_USER,
          password: process.env.LOGGLY_PWD
        }
      });
    }
  }
  return loggerSingleton.prototype.logger;
}

exports = module.exports = loggerSingleton;
