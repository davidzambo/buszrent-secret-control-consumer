package login

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	ClientID     string
	Username     string
	Password     string
	ApiHost      string
	ApiPort      string
	IsAutoLogin  bool
	RefreshToken string
	Slack        struct {
		APIToken    string
		Channel     string
		DevAPIToken string
		DevChannel  string
	}
}

func LoginToWebFlotta(client *http.Client, cfg Config, webFlottaSso string) (statusCode int, err error) {
	log.Println("---loginToWebFlotta")
	callbackUrl := fmt.Sprintf("%s/token", cfg.ApiHost)

	loginURL := webFlottaSso + "/?ClientId=" + cfg.ClientID + "&ApplicationId=Webflotta&CallbackUrl=" + callbackUrl + "&LanguageCode=hu"

	data := url.Values{}
	data.Set("username", cfg.Username)
	data.Set("password", cfg.Password)

	req, _ := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, err
		}
		return resp.StatusCode, fmt.Errorf("loginToWebFlotta error: %s", string(body))
	}
	return resp.StatusCode, nil
}
