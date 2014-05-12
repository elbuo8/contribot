package backends

import (
	"fmt"
	"github.com/elbuo8/contribot/contribot"
	"github.com/sendgrid/sendgrid-go"
)

type Options struct {
	Username string
	Password string
	Alert    []string
	From     string
	Subject  string
}

func Email(opts *Options) contribot.Backend {
	sg := sendgrid.NewSendGridClient(opts.Username, opts.Password)
	mail := sendgrid.NewMail()
	mail.AddTos(opts.Alert)
	mail.SetFrom(opts.From)
	mail.SetSubject(opts.Subject)
	return func(sub *contribot.Submission) {
		sendEmail(mail, sg, sub)
	}
}

func sendEmail(mail *sendgrid.SGMail, sg *sendgrid.SGClient, sub *contribot.Submission) {
	mail.SetText(fmt.Sprintf("New Contributor!\nName: %s\nAddress: %s\nEmail: %s\nSize: %s", sub.Name, sub.Address, sub.Email, sub.Size))
	sg.Send(mail)
}
