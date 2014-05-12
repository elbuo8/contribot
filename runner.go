package main

import (
	"github.com/elbuo8/contribot/backends"
	"github.com/elbuo8/contribot/contribot"
	"log"
	"os"
)

func main() {
	bot := contribot.New()

	bot.Use(func(submission *contribot.Submission) {
		log.Println(submission)
	})
	bot.Use(backends.Email(&backends.Options{
		Username: os.Getenv("SG_USER"),
		Password: os.Getenv("SG_PWD"),
		Alert:    []string{"yamil@sendgrid.com"},
		From:     "yamil@sendgrid.com",
		Subject:  "New Contributor",
	}))
	bot.Run(80)
}
