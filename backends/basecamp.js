var request = require('request');

exports.dispatch = function (userinfo) {
  request({
    method: 'POST',
    url: 'https://basecamp.com/' + process.env.BASECAMP_ACCOUNT +
    '/api/v1/projects/' + process.env.BASECAMP_PROJECT + '/messages.json',
    json: {
      subject: 'Swag for dude/dudette',
      content: userinfo,
      subscribers: process.env.BASECAMP_SUBSCRIBERS.split(',')
    }, auth: {
      user: process.env.BASECAMP_USER,
      pass: process.env.BASECAMP_PWD
    }, headers: {
      'User-Agent': 'ContriBot',
      'Content-Type': 'application/json'
    }
  }, function(e, r, b) {
    console.log(arguments);
  });
};