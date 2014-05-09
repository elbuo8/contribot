package contribot

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func MapServices(m *martini.ClassicMartini) {
	db, err := mgo.Dial(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("%v captured - Closing database connection", sig)
			db.Close()
			os.Exit(1)
		}
	}()

	m.Map(db)
}

func Gandalf(req *http.Request, res http.ResponseWriter, session sessions.Session) {
	if session.Get("user") == "" {
		http.Redirect(res, req, "/auth", http.StatusOK)
	}
}

func MapRoutes(m *martini.ClassicMartini) {
	m.Post("/githook", HandleGitHook)
	m.Get("/auth", AuthGitHub)
	m.Get("/githubAuth", GitHubAuthMiddleware, GetUserFromToken)
	m.Get("/award", Gandalf, AwardUser)
}
