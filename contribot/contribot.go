package contribot

import (
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
	"os"
	"strconv"
)

type contriBot struct {
	Server   *martini.ClassicMartini
	Backends []Backend
}

// New creates contriBot struct
func New() *contriBot {
	err := godotenv.Load() // Make this easier
	if err != nil {
		log.Fatal(err)
	}
	app := martini.Classic()

	mapServices(app)
	mapRoutes(app)
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
	return &contriBot{
		Server: app,
	}
}

func (b *contriBot) Run(port int) {
	b.Server.Map(b.Backends)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), b.Server))
}

func (b *contriBot) Use(backend Backend) {
	b.Backends = append(b.Backends, backend)
}
