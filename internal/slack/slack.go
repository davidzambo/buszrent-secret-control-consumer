package slack

import (
	"github.com/slack-go/slack"
)

func SendMessage(token, channel, msg string) error {
	sl := slack.New(token)
	if _, _, err := sl.PostMessage(channel, slack.MsgOptionText(msg, false)); err != nil {
		return err
	}
	return nil
}

func SendFailedCommunicationWarning(token, channel string) error {
	return SendMessage(token, channel, "Secret Control consumer: nincs érvényes token, a kommunikáció leállt")
}
