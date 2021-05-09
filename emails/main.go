package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"imabad.dev/do/lib/messaging"
	"imabad.dev/do/lib/utils"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

var auth smtp.Auth

func main() {
	utils.LoadConfig()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env, defaulting")
	}
	auth = smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_SERVER"))
	//err = sendEmail("stuart@pomeroys.site", "Hello", "<strong>Hello!</strong>")
	//if err != nil {
	//	log.Fatal("Did not send email", err)
	//} else {
	//	log.Printf("Sent email!")
	//}

	messaging.Setup()
	defer messaging.Close()
	messaging.RegisterQueue(messaging.RegisteredQueue{
		Name: "emails",
		Callback: func(bytes []byte) {
			var request EmailRequest
			json.Unmarshal(bytes, &request)
			log.Printf("Received message %s", bytes)
			err := sendEmail(request.To, request.Subject, request.Content)
			if err != nil {
				log.Printf("Failed to send email %v", err)
			}
		},
	})
	forever := make(chan bool)

	log.Printf("Waiting for messages")
	<-forever
}

func sendEmail(to string, subject string, content string) error {
	message := []byte(fmt.Sprintf("MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\nTo: %s\r\nSubject:%s\r\n\r\n%s", to, subject, content))
	err := smtp.SendMail(os.Getenv("EMAIL_SERVER")+":587", auth, os.Getenv("SEND_AS"), []string{to}, message)
	return err
}
