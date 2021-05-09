package messaging

type OverlayAuthRequest struct {
	Identifier    string `json:"identifier"`
	Pin           string `json:"pin"`
	CorrelationID string `json:"correlationID"`
}

type OverlayAuthResponse struct {
	Identifier    string `json:"identifier"`
	Authenticated bool   `json:"authenticated"`
}

type ChangeStateRequest struct {
	Overlay string      `json:"overlay"`
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
}
