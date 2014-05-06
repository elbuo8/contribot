package main

import (
	"./contribot"
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app := martini.Classic()

	contribot.MapServices(app)
	contribot.MapRoutes(app)

	app.Run()
}
