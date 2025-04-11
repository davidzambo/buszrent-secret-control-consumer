package alerts

import (
	"buszrent-secret-control-consumer/internal/login"
	"buszrent-secret-control-consumer/internal/slack"
	"buszrent-secret-control-consumer/internal/tokens"
	"buszrent-secret-control-consumer/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AlertMessage struct {
	ID            int
	RoadToll      int `json:"Utdij"`
	PlateNumber   string
	DateTime      utils.CustomTime
	MessageType   int
	ObuId         string
	MessageParams []string
	IsRead        bool
	Latitude      float64
	Longitude     float64
	IsResolver    bool
}

type MessageTypes struct {
	Alerts struct {
		Labels map[string]string
	}
	TollroadAlerts map[string]string
}

var messageTypes MessageTypes

var lastFetch = time.Now()

func GetNew(client *http.Client, cfg login.Config, webFlottaUrl string) ([]AlertMessage, error) {
	token := tokens.GetAccessToken()
	if token == "" {
		return nil, errors.New("can't fetch new alert messages, no accessToken")
	}

	alertsUrl := webFlottaUrl + "/alertmessages"

	req, err := http.NewRequest("GET", alertsUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-app-access-token", token)
	req.Header.Set("x-app-client", "Webflotta")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("token is not accepted")
		if err := slack.SendFailedCommunicationWarning(cfg.Slack.DevAPIToken, cfg.Slack.DevChannel); err != nil {
			return nil, fmt.Errorf("slack send error on get new alerts for channel: %s %v", cfg.Slack.DevChannel, err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get alerts request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	var alerts []AlertMessage
	if err := json.Unmarshal(body, &alerts); err != nil {
		return nil, err
	}

	defer func() {
		lastFetch = time.Now()
	}()

	var newAlerts []AlertMessage

	for _, alert := range alerts {
		if alert.DateTime.Time.After(lastFetch) {
			newAlerts = append(newAlerts, alert)
		}
	}

	return newAlerts, nil
}

func FetchMessageTypes(client *http.Client, cfg login.Config, webFlottaUrl, webFlottaSso string, bust int) error {
	typesUrl := webFlottaUrl + "/locales/hu/static.json?bust=" + strconv.Itoa(bust)

	req, err := http.NewRequest("GET", typesUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-app-access-token", tokens.GetAccessToken())
	req.Header.Set("x-app-client", "Webflotta")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("token is not accepted")
		if err := slack.SendFailedCommunicationWarning(cfg.Slack.DevAPIToken, cfg.Slack.DevChannel); err != nil {
			return fmt.Errorf("slack send error on fetch message types to channel: %s %v", cfg.Slack.DevChannel, err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch message types request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var mt MessageTypes
	if err := json.Unmarshal(body, &mt); err != nil {
		return err
	}
	messageTypes = mt
	return nil
}

func JoinAlerts(alerts []AlertMessage) string {
	msg := "Az alábbi webflotta riasztások érkeztek: \n"

	for _, alert := range alerts {
		var alertMsg string
		typeStr := strconv.Itoa(alert.MessageType)

		if alert.RoadToll == 1 {
			alertMsg = messageTypes.TollroadAlerts[typeStr]
		} else {
			alertMsg = messageTypes.Alerts.Labels[typeStr]
		}
		msg += fmt.Sprintf("Rendszám: %s, risztás ideje: %s, üzenet: %s\n",
			alert.PlateNumber,
			alert.DateTime.Time.Local().Format("2006-01-02 15:04"),
			alertMsg)
	}
	return msg
}
