package contribot

import (
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"github.com/martini-contib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
	"os"
)

type ContriBot struct {
	Server   *martini.ClassicMartini
	Backends []Backend
}

func New() *ContriBot {
	err := godotenv.Load() // Make this easier
	if err != nil {
		log.Fatal(err)
	}
	app := martini.Classic()

	MapServices(app)
	MapRoutes(app)
	app.Use(martini.Static("public"))
	app.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET")))
	app.Use(sessions.Sessions("session", store))
	app.Use(csrf.Generate(&csrf.Options{
		Secret:     os.Getenv("CSRF"),
		SessionKey: "user",
		ErrorFunc: func(res http.ResponseWriter) {
			http.Error(res, "CSRF Token Failure", http.StatusUnauthorized)
		},
	}))
	return &ContriBot{
		Server: app,
	}
}

func (b *ContriBot) Run() {
	b.Server.Map(b.Backends)
	b.Server.Run()
}

func (b *ContriBot) Use(backend Backend) {
	b.Backends = append(b.Backends, backend)
}
