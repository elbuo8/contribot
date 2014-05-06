package contribot

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
	"os"
	"os/signal"
)

func MapServices(m *martini.ClassicMartini) {
	session, err := mgo.Dial(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("%v captured - Closing database connection", sig)
			session.Close()
			os.Exit(1)
		}
	}()

	m.Map(session)
}

func MapRoutes(m *martini.ClassicMartini) {
	m.Post("/githook", HandleGitHook)
}
