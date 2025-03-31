package main

import (
	"context"
	"github.com/ardanlabs/conf"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	ClientID         string `conf:"default:--required--"`
	Username         string `conf:"default:--required--"`
	Password         string `conf:"default:--required--"`
	ApiPort          string `conf:"default:8080"`
	SlackApiToken    string `conf:"default:--required--"`
	SlackChannelName string `conf:"default:#slack-bot-messages"`
}

const (
	webFlottaApi = "https://api.webflotta.hu/api"
	webFlottaApp = "https://webflotta.hu/App"
	webFlottaSso = "https://sso.webflotta.hu"
	bust         = 02201446 // TODO: ez az érték a js-ben van környezeti változóként megadva, nem tudtam lekérdezni, és lehet, hogy változni fog valamikor
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found or unable to load.")
	}

	var cfg Config
	if err := conf.Parse(nil, "secret_control_consumer", &cfg); err != nil {
		log.Fatalf("parsing config: %v", err)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	h := Handler{
		client: client,
		config: cfg,
		ssoUrl: webFlottaSso,
	}

	go loginToWebFlotta(cfg, client, webFlottaSso)
	go refreshTokenWorker(cfg.ClientID, webFlottaSso)
	go sendAlertMessageNotificationsWorker(h.client, cfg.SlackApiToken, cfg.SlackChannelName, webFlottaApi, webFlottaApp)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		if err := startHttpServer(&h); err != nil {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-signalChan
	log.Print("Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return nil
}
