
var sendgrid = require('sendgrid')(process.env.SENDGRID_USER, process.env.SENDGRID_PWD);
var logger = require('./../lib/logger');

exports.dispatch = function (userinfo) {
  sendgrid.send({
    to: process.env.EMAIL_TO,
    from: process.env.EMAIL_FROM,
    subject: 'New Contributor!',
    text: JSON.stringify(userinfo)
  }, function (e, r) {
    logger.log(arguments);
  });
};
