package main

import (
	"encoding/json"
	"log"

	"github.com/go-gomail/gomail"
)

func send(data string) (string, error) {

	email := &Email{}
	err := json.Unmarshal([]byte(data), email)
	if err != nil {
		log.Println("Cannot Decode Json of : ", data)
		return "invalidJSON", err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", email.From)
	message.SetHeader("To", email.To)
	message.SetHeader("Subject", email.Subject)
	message.SetBody("text/html", email.Body)

	dialer := gomail.NewDialer(email.Host, email.Port, email.Username, email.Password)

	if err := dialer.DialAndSend(message); err != nil {
		log.Println("Cannot Send Message error : ", err)
		return "cannotSendMessage", err
	}

	return "success", nil
}
