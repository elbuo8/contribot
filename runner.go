package main

import (
	"github.com/elbuo8/contribot/backends"
	"github.com/elbuo8/contribot/contribot"
	"log"
	"os"
	"strings"
)

func main() {
	bot := contribot.New()

	bot.Use(func(submission *contribot.Submission) {
		log.Println(submission)
	})
	bot.Use(backends.Email(&backends.EmailOptions{
		Username: os.Getenv("SG_USER"),
		Password: os.Getenv("SG_PWD"),
		Alert:    []string{"yamil@sendgrid.com"},
		From:     "yamil@sendgrid.com",
		Subject:  "New Contributor",
	}))
	bot.Use(backends.Basecamp(&backends.BasecampOptions{
		Username:    os.Getenv("BC_USER"),
		Password:    os.Getenv("BC_PWD"),
		Subject:     "New Contributor",
		Project:     os.Getenv("BC_PROJECT"),
		Account:     os.Getenv("BC_ACCOUNT"),
		Subscribers: strings.Split(os.Getenv("BC_SUBS"), ","),
	}))
	bot.Run(80)
}
