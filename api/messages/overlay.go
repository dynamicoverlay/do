package messages

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	consts "imabad.dev/do/api/utils"
	"imabad.dev/do/lib/messaging"
	"imabad.dev/do/lib/models"
)

var OverlayQueue = messaging.RegisteredQueue{
	Name: "overlay-auth",
	Callback: func(bytes []byte, d amqp.Delivery) {
		log.Println("Received auth request")
		var request messaging.OverlayAuthRequest
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			log.Print("Failed to decode message", err)
			return
		}
		log.Print(string(bytes))
		var response messaging.OverlayAuthResponse
		if len(request.Identifier) > 0 && len(request.Pin) > 0 {
			log.Println("Verifying auth for", request.Identifier)
			response = VerifyOverlay(request)
		} else {
			response = messaging.OverlayAuthResponse{
				Identifier: request.Identifier, Authenticated: false,
			}
		}
		newBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to encode reply")
			return
		}
		log.Println("Replying to auth request, ", d.ReplyTo, string(newBytes))
		messaging.Reply(d.ReplyTo, request.CorrelationID, "application/json", newBytes)
	},
}

func VerifyOverlay(r messaging.OverlayAuthRequest) messaging.OverlayAuthResponse {
	var overlay models.Overlay
	if err := consts.Db.First(&overlay, "identifier = ? AND pin = ?", r.Identifier, r.Pin).Error; err != nil {
		return messaging.OverlayAuthResponse{Identifier: r.Identifier, Authenticated: false}
	}
	return messaging.OverlayAuthResponse{Identifier: r.Identifier, Authenticated: true}
}
