package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"io/fs"
	"log"

	"path/filepath"

	"github.com/MarselBissengaliyev/ggp-blog/config"
	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.User, data *EmailData, conf *config.Config) {
	to := user.Email
	from := conf.Email_From
	smtpPass := conf.Email_Smtp_Pass
	smtpUser := conf.Email_Smtp_User
	smtpHost := conf.Email_Smtp_Host
	smtpPort := conf.Email_Smtp_Port

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template: ", err)
	}

	if err := template.ExecuteTemplate(&body, "verificationCode.html", &data); err != nil {
		log.Fatal("Could not execute template: ", err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Fatal("could not send email: ", err)
	}
}
