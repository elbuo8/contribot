package main

import (
	"github.com/elbuo8/contribot/contribot"
	"log"
)

func main() {
	bot := contribot.New()

	bot.Use(func(submission *contribot.Submission) {
		log.Println(submission)
	})

	bot.Run("80")
}
