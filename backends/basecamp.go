package backends

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elbuo8/contribot/contribot"
	"net/http"
	"strings"
)

type BasecampOptions struct {
	Username    string
	Password    string
	Account     string
	Project     string
	Subscribers []string
	Subject     string
}

func Basecamp(opts *BasecampOptions) contribot.Backend {
	return func(sub *contribot.Submission) {
		postToBasecamp(opts, sub)
	}
}

func postToBasecamp(opts *BasecampOptions, sub *contribot.Submission) {
	urlStr := fmt.Sprintf("https://basecamp.com/%s/api/v1/projects/%s/messages.json'", opts.Account, opts.Project)
	blob := make(map[string]string)
	blob["subject"] = opts.Subject
	blob["subscribers"] = strings.Join(opts.Subscribers, ",")
	blob["content"] = fmt.Sprintf("New Contributor!\nName: %s\nAddress: %s\nEmail: %s\nSize: %s", sub.Name, sub.Address, sub.Email, sub.Size)
	body, _ := json.Marshal(blob)
	req, _ := http.NewRequest("POST", urlStr, bytes.NewReader(body))
	req.Header.Set("User-Agent", "ContriBot")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(opts.Username, opts.Password)
	http.DefaultClient.Do(req)
}
