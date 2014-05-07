package contribot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	GitHubAPIURL = "https://api.github.com"
	AcceptHeader = "application/vnd.github.v3+json"
	UserAgent    = "ContriBot"
)

func HandleGitHook(req *http.Request, session *mgo.Session) int {
	if req.Header.Get("X-GitHub-Event") != "pull_request" {
		log.Println("Unsed GitHub Payload")
		return http.StatusOK
	}
	log.Println("Received Pull Request Payload")

	err := req.ParseForm()
	if err != nil {
		log.Println(err)
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
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%v", err)
	}
}

func AuthGitHub(req *http.Request, res http.ResponseWriter, params martini.Params, session *mgo.Session, r render.Render) {
	dbSession := session.Copy()
	c := dbSession.DB("contribot").C("contributor")
	status := CheckStatus(c, params["user"])
	if status != 0 { //Auth
		querystring := url.Values{}
		querystring.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
		querystring.Set("redirect_uri", os.Getenv("DOMAIN")+"/award/"+params["user"])
		querystring.Set("scope", "user")
		querystring.Set("state", os.Getenv("SECRET"))
		urlStr := "https://github.com/login/oauth/authorize?" + querystring.Encode()
		log.Println(urlStr)
		http.Redirect(res, req, urlStr, http.StatusFound)
	} else {
		template := make(map[string]string)
		template["message"] = "Sorry, we can't seem to place you."
		template["contactUrl"] = "https://twitter.com/elbuo8" //env
		template["contactValue"] = "@elbuo8"
		r.HTML(http.StatusOK, "error", template)
	}
	dbSession.Close()
}

func AwardUser(req *http.Request, res http.ResponseWriter, params martini.Params, session *mgo.Session) {
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	code := req.Form["code"][0]
	if code == "" {
		http.Redirect(res, req, "/auth/"+params["user"], http.StatusFound)
		return
	}

	dbSession := session.Copy()
	defer dbSession.Close()

	c := dbSession.DB("contribot").C("contributor")
	status := CheckStatus(c, params["user"])
	if status == 0 {
		// Someone is being a troll
		return
	} else if status == 3 {
		//you have been awarded son.
		return
	}
	//move this
	var payload map[string]string
	payload["client_id"] = os.Getenv("GITHUB_CLIENT_ID")
	payload["client_secret"] = os.Getenv("GITHUB_CLIENT_SECRET")
	payload["code"] = code
	body, _ := json.Marshal(payload)
	r, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewReader(body))
	r.Header.Add("Accept", AcceptHeader)
	r.Header.Add("User-Agent", UserAgent)
	ghRes, err := http.DefaultClient.Do(r)

	if err != nil {
		log.Println(err)
		return
	}

	ghPayload, err := ioutil.ReadAll(ghRes.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var ghJSON map[string]interface{}
	err = json.Unmarshal(ghPayload, &ghJSON)
	if err != nil {
		log.Println(err)
		return
	}
	if params["user"] != ghJSON["login"].(string) {
		// someone is being a troll
		return
	}
	// render the form
}
