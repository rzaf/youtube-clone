package email

import (
	"log"
)

func SendVerifcationEmail(userEmail string, username string, link string) error {
	m := make(map[string]any)
	m["username"] = username
	m["link"] = link
	body, err := buildHtml("verification", m)
	if err != nil {
		log.Printf("failed to parse html\n")
		log.Println(err)
		return err
	}
	log.Printf("Sending verification email to %s ..... ", userEmail)
	err = sendEmail([]string{userEmail}, "Verify your email!", body)
	if err != nil {
		return err
	}
	log.Printf("sent\n")
	return nil
}
