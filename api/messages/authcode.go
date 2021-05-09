package messages

type AddAuthCodeRequest struct {
	Code     string `json:"code"`
	TwitchID string `json:"twitchID"`
}