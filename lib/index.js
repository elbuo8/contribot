
var db = require('levelup')('./db');
var GitHubApi = require('github');
var github = new GitHubApi({version: '3.0.0'});
var request = require('request');
var qs = require('querystring');
var backends = process.env.BACKENDS.split(',').map(function(backend) {
  return require(__dirname + '/../backends/' + backend.trim());
});

exports.handleHook = function(req, res) {
  github.pullRequests.get({
    user: req.body.owner,
    repo: req.body.repo,
    number: req.body.number
  }, function (e, pr) {
    if (e !== undefined) {
      console.log(e);
      //logs
    } else {
      if (pr.merged) {
        awardUser(pr.user.login, function (award) {
          if (award) {
            github.authenticate({
              type: 'basic',
              username: process.env.GITHUB_USER,
              password: procces.env.GITHUB_PWD
            });
            github.pullRequests.createCommentReply({
              user: process.env.GITHUB_USER,
              repo: req.body.repo,
              number: req.body.number,
              body: 'Hey dude! Awesome job! We wish to reward you! ' +
              'Please follow the following link. It will ask you to authenticate ' +
              'with your GitHub Account. After that just submit some info and you ' +
              'will be rewarded! \n\nhttp://' + process.env.DOMAIN + '/auth/' + pr.user.login +
              '\n\n Once again, you are AWESOME!',
              in_reply_to: 1
            }, function (e) {
              console.log(e);
            });
          }
        });
      }
    }
  });
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
      github.authenticate({
        type: 'oauth',
        token: b.access_token
      });
      github.user.get({}, function(e, user) {
        if (user.login == req.session.user) {
          res.render('form');
        } else {
          res.send('uh uh, someone messed up');
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
      db.put(user, false);
      cb(true);
    } else if (awarded === 'false') {
      cb(true);
    } else {
      cb(false);
    }
  });
};

awardedUser = function (user, cb) {
  db.put(user, true);
  cb();
};
