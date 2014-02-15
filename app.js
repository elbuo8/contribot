
/**
 * Module dependencies.
 */
require('dotenv').load();
var express = require('express');
var handlers = require('./lib');
var http = require('http');
var path = require('path');

var app = express();

// all environments
app.set('port', process.env.PORT || 3000);
app.set('views', __dirname + '/views');
app.set('view engine', 'jade');
app.use(express.favicon());
app.use(express.cookieParser());
app.use(express.json());
app.use(express.urlencoded());
app.use(express.methodOverride());
app.use(express.session({secret: process.env.SECRET}));
app.use(express.session());
app.use(app.router);
app.use(express.static(path.join(__dirname, 'public')));

app.post('/hook', handlers.handleHook);
app.get('/auth/:user', handlers.auth);
app.get('/award', handlers.award);
app.post('/submit', handlers.submission);

http.createServer(app).listen(app.get('port'), function(){
  console.log('ContriBot listening on port ' + app.get('port'));
});
