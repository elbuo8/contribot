package contribot

import (
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	mapRoutes(app)

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
