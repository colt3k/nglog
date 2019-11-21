package ng

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-mail/mail"
)

//******************* File APPENDER ********************

type MailAppender struct {
	*OutAppender
	smtpServer   string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	smtpFrom     string
	smtpTo       string
	smtpSubject  string
}

func NewMailAppender(filter, smtpServer, smtpUsername, smtpPass, smtpFrom, smtpTo, smtpSubject string, smtpPort int) (*MailAppender, error) {
	if len(smtpServer) == 0 {
		return nil, fmt.Errorf("smtp server required")
	}
	if len(smtpFrom) == 0 {
		return nil, fmt.Errorf("from required")
	}
	if len(smtpTo) == 0 {
		return nil, fmt.Errorf("to required")
	}
	oa := newOutAppender(filter, "")
	t := new(MailAppender)
	t.OutAppender = oa
	t.smtpServer = smtpServer
	t.smtpPort = smtpPort
	t.smtpUsername = smtpUsername
	t.smtpPassword = smtpPass
	t.smtpFrom = smtpFrom
	t.smtpTo = smtpTo
	t.smtpSubject = smtpSubject

	return t, nil
}
func (f *MailAppender) Name() string {
	if len(f.name) > 0 {
		return f.name
	}
	return fmt.Sprintf("%T", f)
}
func (f *MailAppender) DisableColor() bool {
	return f.disableColor
}
func (f *MailAppender) Applicable(msg string) bool {
	if f.filter == "*" {
		return true
	}
	if strings.Index(msg, f.filter) > -1 {
		return true
	}
	return false
}

func (f *MailAppender) Process(msg []byte) {
	m := mail.NewMessage()
	m.SetHeader("From", f.smtpFrom)
	m.SetHeader("To", f.smtpTo)
	m.SetHeader("Subject", f.smtpSubject)
	m.SetBody("text/plain", string(msg))

	d := mail.NewDialer(f.smtpServer, f.smtpPort, f.smtpUsername, f.smtpPassword)
	if f.smtpPort != 25 {
		d.StartTLSPolicy = mail.MandatoryStartTLS
	} else {
		d.SSL = false
		d.StartTLSPolicy = mail.NoStartTLS
	}
	d.Timeout = 5 * time.Second
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("error sending mail\n%+v", err)
	}
}
