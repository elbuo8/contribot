package contribot

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Contributor struct {
	ID         string `bson:"_id"`
	Status     string `bson:"status"`
	StatusCode int    `bson:"status_code"`
}

// Status Legend
// StatusCode - 1 - User is scheduled.
// StatusCode - 2 - User is scheduled and has auth.
// StatusCode - 3 - User has submitted.

func ScheduleContributor(c *mgo.Collection, contributor string) bool {
	var user Contributor
	err := c.FindId(id).One(&existing)
	if err != nil {
		log.Printf("%v", err)
	}
	// Existing should be nill.
	if user != nil {
		return false
	}
	user.ID = contributor
	user.Status = "Scheduled to be rewarded."
	user.StatusCode = 1
	err = c.Insert(user)
	if err != nil {
		log.Printf("%v", err)
	}
	return true
}
