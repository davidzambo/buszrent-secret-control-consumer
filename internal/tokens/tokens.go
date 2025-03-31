package tokens

import (
	"bytes"
	"encoding/json"
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

func FetchNewTokens(refreshToken, clientId, webFlottaSso string) error {
	data := map[string]string{
		"ClientID":     clientId,
		"RefreshToken": refreshToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(webFlottaSso+"/api/refresh", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var refreshTokenResponse struct {
		Data struct {
			AccessToken  string
			RefreshToken string
		}
	}
	if err := json.Unmarshal(body, &refreshTokenResponse); err != nil {
		return err
	}
	SetAccessToken(refreshTokenResponse.Data.AccessToken)
	SetRefreshToken(refreshTokenResponse.Data.RefreshToken)

	return nil
}
