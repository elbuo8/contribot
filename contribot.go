package main

import (
	"./contribot"
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
	"os"
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
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET")))
	app.Use(sessions.Sessions("session", store))
	app.Run()
}
