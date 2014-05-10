# ContriBot

ContriBot is a cool little tool for those whom wish to be grateful for their Open Source contributors. [@SendGrid](http://twitter.com/sendgrid) has several libraries where awesome people submit fixes and we have decided to award them with patches!

![patch](https://pbs.twimg.com/media/BgeS6MwIYAA2I5N.jpg:small )

In order to automate the process of sending stuff, I built a bot which will verify every time a merge has been applied to master. The bot will later collect the users data and post it in all the desired backends (basecamp, trello, lob, email).

## Go-Rewrite

At this time, ContriBot is being re-writen in Go. Most of the core functionality remains the same. Configuration should be simpler. Documentation on how to deploy and contribute to ContriBot will soon follow :)

## Adding it to your repo

Visit the **settings** page in your repo and look for **Hooks and Services**. Use the URL with the following format: http://your-heroku-domain.herokuapp.com/hook

Select only **Pull Requests** to send notifications.