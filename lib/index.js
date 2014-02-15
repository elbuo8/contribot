
var db = require('redis-url').connect(process.env.REDISTOGO_URL);
var request = require('request');
var qs = require('querystring');
var backends = process.env.BACKENDS.split(',').map(function(backend) {
  return require(__dirname + '/../backends/' + backend.trim());
});

var GITHUB_API_URL = 'https://api.github.com/';

exports.handleHook = function(req, res) {
  if (req.headers['x-github-event'] !== 'ping') {
    var pr = req.body;
    var repo = pr.repository.name;
    var owner = pr.repository.owner.login;
    var number = pr.number;
    var contributor = pr.pull_request.user.login;
    var merged = pr.pull_request.merged;
    if (merged) {
      awardUser(contributor, function (award) {
        if (award) {
          request({
            method: 'POST',
            url: GITHUB_API_URL + 'repos/' + owner + '/' + repo + '/issues/' + number +'/comments',
            json: {
              body: 'Hey dude! Awesome job! We wish to reward you! ' +
              'Please follow the following link. It will ask you to authenticate ' +
              'with your GitHub Account. After that just submit some info and you ' +
              'will be rewarded! \n\n' + process.env.DOMAIN + '/auth/' + contributor +
              '\n\n Once again, you are AWESOME!'
            },
            auth: {
              user: process.env.GITHUB_USER,
              pass: process.env.GITHUB_PWD
            }, headers: {
              'User-Agent': 'ContriBot',
              'Content-Type': 'application/json'
            }
          }, function(e, r, b) {
            console.log(arguments);
          });
        }
      });
    }
  }
  res.send(200);
};


exports.auth = function (req, res) {
  awardUser(req.params.user, function(award) {
    if (award) {
      req.session.user = req.params.user;
      var params = qs.stringify({
        client_id: process.env.GITHUB_CLIENT_ID,
        redirect_uri: process.env.DOMAIN + '/award',
        scope: 'user',
        state: process.env.SECRET
      });
      res.redirect('https://github.com/login/oauth/authorize?' + params);
    } else {
      res.send('you have been awarded before brah');
    }
  });
};


exports.award = function (req, res) {
  if (req.query.state === process.env.SECRET) {
    request({url: 'https://github.com/login/oauth/access_token',
      json: {
        client_id: process.env.GITHUB_CLIENT_ID,
        client_secret: process.env.GITHUB_CLIENT_SECRET,
        code: req.query.code
      },
      headers: {Accept: 'application/json'},
      method: 'POST'
    }, function(e, r, b) {
      request({
        method: 'GET',
        url: GITHUB_API_URL + 'user?' + qs.stringify({access_token: b.access_token})
      }, function (e, r, b) {
        if (b.login === req.session.user) {
          res.render('from');
        }
      });
    });
  } else {
    res.send('dont get funny brah');
  }
};

exports.submission = function (req, res) {
  awardUser(req.session.user, function(award) {
    if (award) {
      awardedUser(req.session.user, function () {
        console.log('backendtime');
        backends.forEach(function(backend) {
          backend.dispatch(req.body);
        }); //change to async when not lazy
        res.send('you will be awarded soon :)');
      });
    } else {
      res.send('you have been awarded before brah');
    }
  });
};

awardUser = function (user, cb) {
  /*
    cb return params
    true - if user is not db and should be inserted with false
    false - if user is on db (regardless of status)
  */
  db.get(user, function (e, awarded) {
    if (e) {
      cb(false);
    } else if (awarded === null ||awarded === 'false') {
      db.set(user, false);
      cb(true);
    } else {
      cb(false);
    }
  });
};

awardedUser = function (user, cb) {
  db.set(user, true);
  cb();
};
