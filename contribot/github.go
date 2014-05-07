package contribot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	GitHubAPIURL = "https://api.github.com"
	AcceptHeader = "application/vnd.github.v3+json"
)

func HandleGitHook(req *http.Request, session *mgo.Session) int {
	if req.Header.Get("X-GitHub-Event") != "pull_request" {
		log.Println("Unsed GitHub Payload")
		return http.StatusOK
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
		c := dbSession.DB("contribot").C("contributor")
		userInfo := pullRequest["user"].(map[string]interface{})
		scheduled := ScheduleContributor(c, userInfo["login"].(string))
		if scheduled {
			log.Printf("New Contributor: %s", userInfo["login"])
			repository := payload["repository"].(map[string]interface{})
			repoName := repository["full_name"].(string)
			pullRequestNumber := fmt.Sprintf("%.0f", pullRequest["number"].(float64))
			go PostRewardInvite(repoName, pullRequestNumber)
		}
		// Clean up
		dbSession.Close()
	}

	return http.StatusOK
}

func PostRewardInvite(repoName, prNumber string) {
	requestUrl := GitHubAPIURL + "/repos/" + repoName + "/issues/" + prNumber + "/comments"
	payload := make(map[string]string)
	payload["body"] = "Hey! Awesome job! We wish to reward you! " +
		"Please follow the following link. It will ask you to authenticate " +
		"with your GitHub Account. After that just submit some info and you " +
		"will be rewarded! \n\n" + "[Click Here!](" + os.Getenv("DOMAIN") + "/auth/" + repoName[0:strings.Index(repoName, "/")] + ")" +
		"\n\n Once again, you are AWESOME!"
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", requestUrl, bytes.NewReader(body))
	req.Header.Add("Accept", AcceptHeader)
	req.Header.Add("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%v", err)
	}
}
