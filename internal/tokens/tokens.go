package tokens

import (
	"buszrent-secret-control-consumer/internal/slack"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var accessToken string
var refreshToken string

func GetAccessToken() string {
	return accessToken
}
func GetRefreshToken() string {
	return refreshToken
}

func SetAccessToken(token string) {
	accessToken = token
}

func SetRefreshToken(token string) {
	refreshToken = token
}

func FetchNewTokens(slackDevToken, slackDevChannel, refreshToken, clientId, webFlottaSso string) (statusCode int, err error) {
	if refreshToken == "" {
		return 400, errors.New("can't fetch new tokens, no refreshToken")
	}

	data := map[string]string{
		"ClientID":     clientId,
		"RefreshToken": refreshToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(webFlottaSso+"/api/refresh", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var refreshTokenResponse struct {
		StatusCode int
		Data       struct {
			AccessToken  string
			RefreshToken string
		}
		Errors interface{}
	}

	if err := json.Unmarshal(body, &refreshTokenResponse); err != nil {
		return 0, err
	}

	if refreshTokenResponse.StatusCode != 200 {
		SetRefreshToken("")
		SetAccessToken("")
		if err := slack.SendFailedCommunicationWarning(slackDevToken, slackDevChannel); err != nil {
			return 0, fmt.Errorf("slack send error on fetch new token to channel: %s %v", slackDevChannel, err)
		}
		return refreshTokenResponse.StatusCode, fmt.Errorf("%v", refreshTokenResponse.Errors)
	}

	SetAccessToken(refreshTokenResponse.Data.AccessToken)
	SetRefreshToken(refreshTokenResponse.Data.RefreshToken)

	return resp.StatusCode, nil
}
