package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	WaitGroup   *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

type Message struct {
	From          string
	FromName      string
	To            string
	Subject       string
	Attachments   []string
	AttachmentMap map[string]string
	Data          any
	DataMap       map[string]any
	Template      string
}

// A function who is listener on the MailerChan

func (app *Config) listenForMail() {

	for {
		select {
		case msg := <-app.Mailer.MailerChan:
			go app.Mailer.sendMail(msg, app.Mailer.ErrorChan)
		case err := <-app.Mailer.ErrorChan:
			app.ErrorLog.Println(err) //You can handle accordingly
		case <-app.Mailer.DoneChan:
			return
		}
	}

}

func (m *Mail) sendMail(msg Message, errorChan chan error) {

	defer m.WaitGroup.Done()

	if msg.Template == "" {
		msg.Template = "mail"
	}

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	if msg.AttachmentMap == nil {
		msg.AttachmentMap = make(map[string]string)
	}

	if len(msg.DataMap) == 0 {
		msg.DataMap = make(map[string]any)
	}
	msg.DataMap["message"] = msg.Data

	//bulid html mail
	formattedMessage, err := m.bulidHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	//bulid plain text mail
	plainMessage, err := m.bulidPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// We have to create a mail server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncrption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, val := range msg.Attachments {
			email.AddAttachment(val)
		}
	}

	if len(msg.AttachmentMap) > 0 {
		for key, val := range msg.AttachmentMap {
			email.AddAttachment(val, key)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}
}

func (m *Mail) bulidHTMLMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", nil
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", nil
	}

	formattedMessageReady := tpl.String()
	formattedMessageReady, err = m.inlineCSS(formattedMessageReady)
	if err != nil {
		return "", nil
	}

	return formattedMessageReady, nil
}

func (m *Mail) bulidPlainTextMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", nil
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", nil
	}

	plainMessage := tpl.String()
	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {

	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncrption(e string) mail.Encryption {

	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
