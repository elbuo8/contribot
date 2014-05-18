package backends

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elbuo8/contribot/contribot"
	"net/http"
)

type HipChatOptions struct {
	Token  string
	Room   string
	Color  string
}

func HipChat(opts *HipChatOptions) contribot.Backend {
	return func(sub *contribot.Submission) {
		sendNotification(opts, sub)
	}
}

func sendNotification(opts *HipChatOptions, sub *contribot.Submission) {
	urlStr := fmt.Sprintf("https://api.hipchat.com/v2/room/%s/notification?auth_token=%s", opts.Room, opts.Token)
	blob := make(map[string]string)
	blob["color"] = opts.Color
  blob["message_format"] = "html"
  blob["message"] = fmt.Sprintf("New Contributor!<br/><strong>Name:</strong> %s<br/><strong>Address:</strong> %s<br/><strong>Email:</strong> %s<br/><strong>Size:</strong> %s", sub.Name, sub.Address, sub.Email, sub.Size)
	body, _ := json.Marshal(blob)
	req, _ := http.NewRequest("POST", urlStr, bytes.NewReader(body))
	req.Header.Set("User-Agent", "ContriBot")
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}
