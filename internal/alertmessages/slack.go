package alertmessages

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"strconv"
	"time"
)

type SlackMessage struct {
	PlateNumber string
	DateTime    time.Time
	Type        string
}

func SendSlackNotifications(sl *slack.Client, slackChannelName string, alerts []AlertMessage) error {
	msg := "Az alábbi webflotta riasztások érkeztek: \n"

	msgList := GetMessageTypes()
	if len(msgList.TollroadAlerts) == 0 {
		return errors.New("empty alert type messages")
	}

	for _, alert := range alerts {
		var alertMsg string
		typeStr := strconv.Itoa(alert.MessageType)

		if alert.RoadToll == 1 {
			alertMsg = msgList.TollroadAlerts[typeStr]
		} else {
			alertMsg = msgList.Alerts.Labels[typeStr]
		}
		msg += fmt.Sprintf("Rendszám: %s, risztás ideje: %s, üzenet: %s\n",
			alert.PlateNumber,
			alert.DateTime.Time.Local().Format("2006-01-02 15:04"),
			alertMsg)
	}

	_, _, err := sl.PostMessage(slackChannelName, slack.MsgOptionText(msg, false))
	if err != nil {
		return err
	}
	return nil
}
