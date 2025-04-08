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
	refreshTokenTimer  = time.Second * 30
	alertMessagesTimer = time.Minute * 1
)

func refreshTokenWorker(slackDevToken, slackDevChannel, clientId string) {
	log.Println("---start refreshTokenWorker")
	for {
		refreshToken := tokens.GetRefreshToken()
		if refreshToken == "" {
			log.Println("unknown refreshToken")
			time.Sleep(refreshTokenTimer)
			continue
		}

		log.Printf("fetch tokens with refreshToken %s\n", refreshToken)
		if statusCode, err := tokens.FetchNewTokens(slackDevToken, slackDevChannel, refreshToken, clientId, webFlottaSso); err != nil {
			log.Printf("Fetch token status: %d, error: %v\nrefreshToken: %s\nclientId: %s\n", statusCode, err, refreshToken, clientId)
		}
		time.Sleep(refreshTokenTimer)
	}
}

func sendAlertMessageNotificationsWorker(client *http.Client, cfg Config) {
	log.Println("---start sendAlertMessageNotificationsWorker")
	for {
		time.Sleep(time.Second * 5)

		if tokens.GetAccessToken() == "" {
			log.Println("unknown accessToken")
			continue
		}

		log.Println("---fetching message types")

		if err := alerts.FetchMessageTypes(client, login.Config(cfg), webFlottaApp, webFlottaSso, bust); err != nil {
			log.Printf("Error on fetching alert type messages: %v\n", err)
			if err := slack.SendMessage(cfg.Slack.DevAPIToken, cfg.Slack.DevChannel, "Secret Control consumer: Error on fetching alert type messages"); err != nil {
				log.Printf("Failed to send message to message type error to channel: %s %v\n", cfg.Slack.DevChannel, err)
			}
			return
		}

		log.Println("---start fetching new alerts")

		for {
			alertList, err := alerts.GetNew(client, login.Config(cfg), webFlottaApi, webFlottaSso)
			if err != nil {
				log.Println(err)
			}

			if alertList != nil {
				jointAlerts := alerts.JoinAlerts(alertList)

				log.Println(jointAlerts)

				if err := slack.SendMessage(cfg.Slack.APIToken, cfg.Slack.Channel, jointAlerts); err != nil {
					log.Printf("Failed to send message to alert message to channel: %s %v\n", cfg.Slack.Channel, err)
				}
			}

			time.Sleep(alertMessagesTimer)
		}
	}
}
