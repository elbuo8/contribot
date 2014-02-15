
var sendgrid = require('sendgrid')(process.env.SENDGRID_USER, process.env.SENDGRID_PWD);

exports.dispatch = function (userinfo) {
  console.log(userinfo);
  sendgrid.send({
    to: process.env.EMAIL_TO,
    from: process.env.EMAIL_FROM,
    subject: 'New Contributor!',
    text: userinfo
  }, function (e, r) {
    console.log(arguments);
  });
};


