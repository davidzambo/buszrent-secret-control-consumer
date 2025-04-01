package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func loginToWebFlotta(cfg Config, client *http.Client, webFlottaSso string) {
	callbackUrl := "http://localhost:" + cfg.ApiPort + "/get_tokens" // TODO: itt jó lesz a localhost?

	loginURL := webFlottaSso + "/?ClientId=" + cfg.ClientID + "&ApplicationId=Webflotta&CallbackUrl=" + callbackUrl + "&LanguageCode=hu"

	data := url.Values{}
	data.Set("username", cfg.Username)
	data.Set("password", cfg.Password)

	req, _ := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("loginToWebFlotta: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("login request failed with status: %d", resp.StatusCode)
	}
}
