package main

import (
	"buszrent-secret-control-consumer/internal/alertmessages"
	"buszrent-secret-control-consumer/internal/tokens"
	"github.com/slack-go/slack"

	"log"
	"net/http"
	"time"
)

const (
	refreshTokenTimer  = time.Minute * 1
	alertMessagesTimer = time.Minute * 1
)

func refreshTokenWorker(clientId, webFlottaSso string) {
	for {
		refreshToken := tokens.GetRefreshToken()
		if refreshToken != "" {
			if err := tokens.FetchNewTokens(refreshToken, clientId, webFlottaSso); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(refreshTokenTimer)
	}
}

func sendAlertMessageNotificationsWorker(client *http.Client, slackApiToken, slackChannel, webFlottaApi, webFlottaApp string) {
	for {
		if tokens.GetAccessToken() == "" {
			time.Sleep(time.Second * 5)
		} else {
			if err := alertmessages.FetchMessageTypes(client, webFlottaApp, bust); err != nil {
				log.Fatalf("fetching alert type messages: %v", err)
			}

			for {
				alerts, err := alertmessages.GetNewAlerts(client, webFlottaApi)
				if err != nil {
					log.Fatal(err)
				}

				if alerts != nil {
					slackApi := slack.New(slackApiToken)
					if err := alertmessages.SendSlackNotifications(slackApi, slackChannel, alerts); err != nil {
						log.Fatal(err)
					}
				}

				time.Sleep(alertMessagesTimer)
			}
		}
	}
}
