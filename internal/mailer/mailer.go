package mailer

import (
	"activator/internal"
	"activator/internal/app/models"
	"activator/internal/config"
	"bytes"
	"html/template"
	"log/slog"
	"net/url"
	"time"

	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	cfg    config.Config
	dialer *mail.Dialer
	sender string
	logger *slog.Logger
}

var templateMessage = `
Hello, {{.RecipientName}}!
Your activation link: {{.Link}}
`

func New(cfg config.Config, logger *slog.Logger, host string, port int, sender string) Mailer {
	dialer := mail.Dialer{Host: host, Port: port}
	dialer.Timeout = 5 * time.Second

	return Mailer{
		cfg:    cfg,
		dialer: &dialer,
		sender: sender,
		logger: logger,
	}
}

func (m Mailer) Send(recipient models.User, token string) error {
	templ, err := template.New("email").Parse(templateMessage)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "template parsing error")
	}

	type LetterData struct {
		RecipientName string
		Link          string
	}

	var messageData LetterData
	messageData.RecipientName = recipient.Name
	messageData.Link = makeActivationUri(m.cfg, token)
	m.logger.Debug("activation link", "ref", messageData.Link)

	plainBody := new(bytes.Buffer)
	err = templ.ExecuteTemplate(plainBody, "email", messageData)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "template execution error")
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient.Email)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", "Subscription activation letter")
	msg.SetBody("text/plain", plainBody.String())

	m.logger.Debug("message body", "msg", plainBody.String())

	attempts := 2
	for i := 0; i < attempts; i++ {
		err = m.dialer.DialAndSend(msg)
		if err == nil {
			// success
			m.logger.Info("Successfully sent email to:", "recipient", recipient.Email)
			return nil
		} else {
			m.logger.Error("DialAndSend error", "err", err.Error())
		}

		time.Sleep(100 * time.Millisecond)
	}

	return internal.NewErrorf(internal.ErrorCodeUnknown, "%d attempts, message was not sent", attempts)
}

func makeActivationUri(cfg config.Config, token string) string {
	q := make(url.Values)
	q.Add("token", token)
	url := url.URL{
		Scheme:   "http",
		Host:     cfg.ServerAddr,
		Path:     cfg.ActivationPath,
		RawQuery: q.Encode(),
	}

	return url.String()
}
