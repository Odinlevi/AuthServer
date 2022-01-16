package sendsmtp

import (
	log "github.com/sirupsen/logrus"
	"net/smtp"
)

type ISender interface {
	Send(to string, subjectMsg string, bodyMsg string)
}

type sender struct {
	from string
	pass string
	host string
	port string

	auth smtp.Auth
}

func NewSender(from string, pass string, host string, port string) ISender {
	data := &sender{
		from: from,
		pass: pass,
		host: host,
		port: port,

		auth: smtp.PlainAuth("", from, pass, host),
	}

	return data
}

func (s *sender) Send(to string, subjectMsg string, bodyMsg string) {
	body := "From: " + s.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subjectMsg + "\n" +
		bodyMsg

	err := smtp.SendMail(s.host+":"+s.port,
		s.auth,
		s.from,
		[]string{to},
		[]byte(body))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
