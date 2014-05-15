package contribot

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func mapServices(m *martini.ClassicMartini) {
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

func gandalf(req *http.Request, res http.ResponseWriter, session sessions.Session) {
	if session.Get("user") == "" {
		http.Redirect(res, req, "/auth", http.StatusFound)
		return
	}
}

func mapRoutes(m *martini.ClassicMartini) {
	m.Post("/githook", handleGitHook)
	m.Get("/auth", authGitHub)
	m.Get("/githubAuth", gitHubAuthMiddleware, getUserFromToken)
	m.Get("/award", gandalf, awardUser)
	m.Post("/submission", gandalf, csrf.Validate, handleSubmission)
}
