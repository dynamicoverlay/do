package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TwitchUser struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Email       string `json:"email"`
}

var httpClient = &http.Client{}

func GetUser(accessToken string, clientID string) (*TwitchUser, error) {
	if len(accessToken) < 1 {
		return nil, fmt.Errorf("invalid access token")
	}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Client-ID", clientID)
	req.Header.Add("Authorization", fmt.Sprintf(`Bearer %s`, accessToken))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string][]TwitchUser
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	users := response["data"]
	if len(users) > 1 {
		return nil, fmt.Errorf("more than one user returned")
	}
	return &users[0], nil
}
