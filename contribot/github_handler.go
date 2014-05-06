package contribot

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

func HandleGitHook(req *http.Request, res http.ResponseWriter, session *mgo.Session) {
	if req.Header.Get("X-GitHub-Event") != "pull_request" {
		log.Println("Unsed GitHub Payload")
		res.WriteHeader(http.StatusOK) // Exit quickly
	}
	log.Println("Received Pull Request Payload")

	err := req.ParseForm()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	rawPayload := req.PostForm.Get("payload")
	var payload map[string]interface{}
	err = json.Unmarshal([]byte(rawPayload), &payload)

	pullRequest := payload["pull_request"].(map[string]interface{})
	mergedPullRequest := pullRequest["merged"].(bool)

	if mergedPullRequest {
		dbSession := session.Copy()
		c, err = dbSession.DB("contribot").C("contributor")
		if err != nil {
			log.Printf("%v", err)
		}
		userInfo := pullRequest["user"].(map[string]interface{})
		scheduled := ScheduleContributor(c, userInfo["login"])
		if scheduled {
			log.Printf("New Contributor: %s", userInfo["login"])
		}
		// Clean up
		dbSession.Close()
	}

	res.WriteHeader(http.StatusOK)
}
