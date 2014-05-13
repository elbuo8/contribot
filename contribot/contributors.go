package contribot

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type userT struct {
	ID         string `bson:"_id"`
	StatusCode int    `bson:"status_code"`
}

// Status Legend
// StatusCode - 1 - User is scheduled.
// StatusCode - 2 - User is scheduled and has auth.
// StatusCode - 3 - User has submitted.

func scheduleContributor(db *mgo.Database, contributor string) bool {
	var user userT
	c := db.C("contributor")
	err := c.FindId(contributor).One(&user)
	if err != mgo.ErrNotFound {
		return false // User shouldn't be in DB
	}
	user.ID = contributor
	user.StatusCode = 1
	err = c.Insert(user)
	if err != nil {
		log.Printf("%v", err)
	}
	return true
}

func checkStatus(db *mgo.Database, contributor string) int {
	var user userT
	err := db.C("contributor").FindId(contributor).One(&user)
	if err != nil {
		return 0
	}
	return user.StatusCode
}

func userHasAuth(db *mgo.Database, contributor string) error {
	return db.C("contributor").UpdateId(contributor, bson.M{"$set": bson.M{"status_code": 2}})
}

func userHasSubmitted(db *mgo.Database, contributor string) error {
	return db.C("contributor").UpdateId(contributor, bson.M{"$set": bson.M{"status_code": 3}})
}
