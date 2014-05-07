package contribot

import (
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
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
	err := c.FindId(contributor).One(&user)
	if err != mgo.ErrNotFound {
		return false // User shouldn't be in DB
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

func CheckStatus(c *mgo.Collection, contributor string) int {
	var user Contributor
	err := c.FindId(contributor).One(&user)
	if err != nil {
		return 0
	}
	return user.StatusCode
}
