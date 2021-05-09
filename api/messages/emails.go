package messages

import (
	"encoding/json"

	"imabad.dev/do/lib/messaging"
)

var EmailQueue = messaging.RegisteredQueue{
	Name: "emails",
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func SendEmail(request EmailRequest) error {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	messaging.Publish(EmailQueue, "application/json", bodyBytes)
	return nil
}
