package email

import (
	"fmt"
	"log"
	"net/smtp"
	"youtube-clone/database/helpers"
)

var (
	smtpAddr  string
	smtpAuth  smtp.Auth
	fromEmail = "john@doe.com"
)

func InitVars() {
	host := helpers.FatalIfEmptyVar("SMTP_HOST")
	port := helpers.FatalIfEmptyVar("SMTP_PORT")
	smtpAddr = fmt.Sprintf("%s:%s", host, port)

	username := helpers.FatalIfEmptyVar("SMTP_USERNAME")
	secret := helpers.FatalIfEmptyVar("SMTP_SECRET")

	smtpAuth = smtp.CRAMMD5Auth(username, secret)
}

func sendEmail(to []string, subjectStr string, htmlBody string) error {
	subject := "Subject: " + subjectStr + "\n"
	from := "From: " + fromEmail + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + from + mime + htmlBody)

	err := smtp.SendMail(smtpAddr, smtpAuth, fromEmail, to, msg)
	if err != nil {
		log.Println("Failed to send email!!!")
		log.Println(err)
		return err
	}
	return nil
}
