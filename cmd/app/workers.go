package main

import (
	"buszrent-secret-control-consumer/internal/alerts"
	"buszrent-secret-control-consumer/internal/login"
	"buszrent-secret-control-consumer/internal/slack"
	"buszrent-secret-control-consumer/internal/tokens"
	"log"
	"net/http"
	"time"
)

const (
	refreshTokenTimer  = time.Minute * 1
	alertMessagesTimer = time.Minute * 1
)

func refreshTokenWorker(slackDevToken, slackDevChannel, clientId string) {
	for {
		log.Println("---start refreshTokenWorker")
		refreshToken := tokens.GetRefreshToken()
		if refreshToken != "" {
			log.Println("fetch tokens")
			if statusCode, err := tokens.FetchNewTokens(slackDevToken, slackDevChannel, refreshToken, clientId, webFlottaSso); err != nil {
				log.Printf("Fetch token status: %d, error: %v\n", statusCode, err)

			}
		}
		time.Sleep(refreshTokenTimer)
	}
}

func sendAlertMessageNotificationsWorker(client *http.Client, cfg Config) {
	for {
		time.Sleep(time.Second * 5) // TODO ez így jó?
		log.Println("---start sendAlertMessageNotificationsWorker")

		if tokens.GetAccessToken() != "" {
			log.Println("---start fetching new alerts")
			if err := alerts.FetchMessageTypes(client, login.Config(cfg), webFlottaApp, webFlottaSso, bust); err != nil {
				log.Printf("fetching alert type messages: %v\n", err)
			}

			for {
				alertList, err := alerts.GetNew(client, login.Config(cfg), webFlottaApi, webFlottaSso)
				if err != nil {
					log.Println(err)
				}

				if alertList != nil {
					if err := slack.SendMessage(cfg.Slack.APIToken, cfg.Slack.Channel, alerts.CreateSlackMessage(alertList)); err != nil {
						log.Println(err)
					}
				}

				time.Sleep(alertMessagesTimer)
			}
		}
	}
}
