package main

import (
	"./contribot"
	"log"
)

func main() {
	bot := contribot.New()
	/*
		bot.Use(func(submission *contribot.Submission) {
	    log.Println(submission)
		})
	*/
	bot.Run()
}
