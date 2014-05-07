package main

import (
	"./contribot"
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/render"
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
	app.Use(martini.Static("public"))
	app.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	app.Run()
}
