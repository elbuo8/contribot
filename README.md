# ContriBot

ContriBot is a cool little tool for those whom wish to be grateful for their Open Source contributors. [@SendGrid](http://twitter.com/sendgrid) has several libraries where awesome people submit fixes and we have decided to award them with patches!

![patch](https://pbs.twimg.com/media/BgeS6MwIYAA2I5N.jpg:small )

In order to automate the process of sending stuff, I built a bot which will verify every time a merge has been applied to master. The bot will later collect the users data and post it in all the desired backends (basecamp, trello, lob, email).

## Running your own ContriBot

```bash
$ git clone https://github.com/elbuo8/contribot.git
$ cd contribot && npm install
$ mv sample.env .env
```

Modify **.env** with the required values. The rest are optional depending on your backend selection.

* BACKENDS #Comma separated list of backends
* SECRET #Random string
* DOMAIN
* GITHUB\_CLIENT\_ID
* GITHUB\_CLIENT\_SECRET
* GITHUB\_USER
* GITHUB\_PWD
* LOGGLY\_USER
* LOGGLY\_TOKEN
* LOGGLY\_PWD
* LOGGLY\_SUBDOMAIN

```bash
$ heroku create _____
$ heroku addons:add redistogo
$ heroku config:push
```

After all this, specify your backends and push to Heroku.

## Adding it to your repo

Visit the **settings** page in your repo and look for **Hooks and Services**. Use the URL with the following format: http://your-heroku-domain.herokuapp.com/hook

Select only **Pull Requests** to send notifications.
