package contribot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	gitHubAPIURL = "https://api.github.com"
	acceptHeader = "application/json"
	userAgent    = "ContriBot"
)

func handleGitHook(req *http.Request, res http.ResponseWriter, db *mgo.Session) {
	if req.Header.Get("X-GitHub-Event") != "pull_request" {
		log.Println("Unsed GitHub Payload")
		res.WriteHeader(http.StatusOK)
		return
	}

	log.Println("Received Pull Request Payload")

	var rawPayload []byte
	var err error
	if req.Header.Get("content-type") == "application/x-www-form-urlencoded" {
		err = req.ParseForm()
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		payload := req.PostForm.Get("payload")
		rawPayload = []byte(payload)
	} else {
		rawPayload, err = ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var payload map[string]interface{}
	err = json.Unmarshal(rawPayload, &payload)

	pullRequest := payload["pull_request"].(map[string]interface{})
	mergedPullRequest := pullRequest["merged"].(bool)
	userInfo := pullRequest["user"].(map[string]interface{})
	repository := payload["repository"].(map[string]interface{})
	repoName := repository["full_name"].(string)
	username := userInfo["login"].(string)
	filter, err := isCollaborator(repoName, username)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if mergedPullRequest && !filter {
		dbSession := db.Copy()
		c := dbSession.DB("contribot").C("contributor")
		scheduled := scheduleContributor(c, username)
		if scheduled {
			log.Printf("New Contributor: %s", username)
			pullRequestNumber := fmt.Sprintf("%.0f", pullRequest["number"].(float64))
			go postRewardInvite(repoName, pullRequestNumber)
		}
		// Clean up
		dbSession.Close()
	}
	res.WriteHeader(http.StatusOK)
}

func isCollaborator(repoName, username string) (bool, error) {
	res, err := http.Get(fmt.Sprintf(gitHubAPIURL+"/repos/%s/collaborators/%s", repoName, username))
	if err != nil {
		return false, err
	}
	return res.StatusCode == http.StatusNoContent, err
}

func postRewardInvite(repoName, prNumber string) {
	requestURL := gitHubAPIURL + "/repos/" + repoName + "/issues/" + prNumber + "/comments"
	payload := make(map[string]string)
	payload["body"] = "Thanks for contributing to SendGrid Open Source! " +
		"We think it's awesome when community members contribute " +
		"to our projects and want to celebrate that." +
		"\n\nThe following link will ask you to authenticate with Github " +
		"(so we can verify your identity), it will then ask for a little bit more " +
		"information and we'll then send you a thanks for contributing." +
		"\n\n**[Click Here to Continue »](" + os.Getenv("DOMAIN") + "/auth)**" +
		"\n\nOnce again, thank you!"
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", requestURL, bytes.NewReader(body))
	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
}

func authGitHub(req *http.Request, res http.ResponseWriter) {
	querystring := url.Values{}
	querystring.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	querystring.Set("redirect_uri", os.Getenv("DOMAIN")+"/githubAuth")
	querystring.Set("scope", "user")
	urlStr := "https://github.com/login/oauth/authorize?" + querystring.Encode()
	http.Redirect(res, req, urlStr, http.StatusFound)
}

func gitHubAuthMiddleware(req *http.Request, res http.ResponseWriter, r render.Render, c martini.Context) {
	// Verify origin is GH
	template := make(map[string]string)
	template["contactUrl"] = os.Getenv("CONTACT_URL")
	template["contactValue"] = os.Getenv("CONTACT_VALUE")
	template["message"] = "There was an authenticating your account."
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusBadRequest, "error", template)
		return
	}
	if len(req.Form["code"]) != 1 {
		r.HTML(http.StatusUnauthorized, "error", template)
		return
	}
	// If legit, attempt to get token
	payload := make(map[string]string)
	payload["client_id"] = os.Getenv("GITHUB_CLIENT_ID")
	payload["client_secret"] = os.Getenv("GITHUB_CLIENT_SECRET")
	payload["code"] = req.Form["code"][0]
	body, _ := json.Marshal(payload)
	ghReq, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewReader(body))
	ghReq.Header.Add("Content-Type", acceptHeader)
	ghReq.Header.Add("Accept", acceptHeader)
	ghReq.Header.Add("User-Agent", userAgent)
	ghRes, err := http.DefaultClient.Do(ghReq)

	// check status code
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusServiceUnavailable, "error", template)
		return
	}
	ghPayload, err := ioutil.ReadAll(ghRes.Body)
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusInternalServerError, "error", template)
		return
	}
	var ghJSON map[string]interface{}
	err = json.Unmarshal(ghPayload, &ghJSON)
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusInternalServerError, "error", template)
		return
	}
	token, ok := ghJSON["access_token"].(string)
	if !ok {
		r.HTML(http.StatusOK, "error", template)
		return
	}
	c.Map(token)
	c.Next()
	http.Redirect(res, req, "/award", http.StatusFound)
}

func getUserFromToken(db *mgo.Session, r render.Render, token string, session sessions.Session) {
	template := make(map[string]string)
	template["contactUrl"] = os.Getenv("CONTACT_URL")
	template["contactValue"] = os.Getenv("CONTACT_VALUE")
	template["message"] = "GitHub seems to have troubles :/"

	qs := url.Values{}
	qs.Set("access_token", token)
	ghReq, _ := http.NewRequest("GET", gitHubAPIURL+"/user?"+qs.Encode(), nil)
	ghReq.Header.Add("User-Agent", userAgent)
	ghRes, err := http.DefaultClient.Do(ghReq)
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusServiceUnavailable, "error", template)
		return
	}
	ghPayload, err := ioutil.ReadAll(ghRes.Body)
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusInternalServerError, "error", template)
		return
	}
	ghRes.Body.Close()
	var ghJSON map[string]interface{}
	err = json.Unmarshal(ghPayload, &ghJSON)
	if err != nil {
		log.Println(err)
		r.HTML(http.StatusInternalServerError, "error", template)
		return
	}

	user, ok := ghJSON["login"].(string)
	if !ok {
		log.Println("Obtaining username from request failed.")
		r.HTML(http.StatusInternalServerError, "error", template)
	}
	session.Set("user", user)
}

func awardUser(db *mgo.Session, session sessions.Session, r render.Render, x csrf.CSRF) {
	template := make(map[string]string)
	template["contactUrl"] = os.Getenv("CONTACT_URL")
	template["contactValue"] = os.Getenv("CONTACT_VALUE")
	dbSession := db.Copy()
	user := session.Get("user").(string)
	status := checkStatus(dbSession.DB("contribot").C("contributor"), user)
	if status == 0 {
		template["message"] = "Can't seem to find records of you :/"
		r.HTML(http.StatusUnauthorized, "error", template)
	} else if status == 1 {
		err := userHasAuth(dbSession.DB("contribot").C("contributor"), user)
		if err != nil {
			log.Println(err)
			template["message"] = "Uh oh! Please report this :("
			r.HTML(http.StatusInternalServerError, "error", template)
		} else {
			r.HTML(http.StatusOK, "form", x.GetToken())
		}
	} else if status == 2 {
		r.HTML(http.StatusOK, "form", x.GetToken())
	} else if status == 3 {
		template["message"] = "Hey buddy, it seems you have been awarded before."
		r.HTML(http.StatusUnauthorized, "error", template)
	}
	dbSession.Close()
}

func handleSubmission(req *http.Request, r render.Render, db *mgo.Session, session sessions.Session, backends []Backend) {
	template := make(map[string]string)
	template["contactUrl"] = os.Getenv("CONTACT_URL")
	template["contactValue"] = os.Getenv("CONTACT_VALUE")
	template["message"] = "Something went wrong :'("
	err := req.ParseForm()
	if err != nil {
		r.HTML(http.StatusBadRequest, "error", template)
	}
	user := session.Get("user").(string)
	dbSession := db.Copy()
	err = userHasSubmitted(dbSession.DB("contribot").C("contributor"), user)

	if err != nil {
		log.Println(err)
		r.HTML(http.StatusInternalServerError, "error", template)
	} else {
		submission := &Submission{
			Name:    req.PostForm.Get("name"),
			Address: req.PostForm.Get("address"),
			Email:   req.PostForm.Get("email"),
			Size:    req.PostForm.Get("size"),
		}
		for i := 0; i < len(backends); i++ {
			go backends[i](submission)
		}
		r.HTML(http.StatusOK, "success", nil)
	}
}
