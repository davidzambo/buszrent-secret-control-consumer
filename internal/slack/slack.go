package slack

import (
	"github.com/slack-go/slack"
)

func GetTokenErrorMessage() string {
	return "Secret control consumer: hibás token"
}

func SendMessage(token, channel, msg string) error {
	sl := slack.New(token)
	if _, _, err := sl.PostMessage(channel, slack.MsgOptionText(msg, false)); err != nil {
		return err
	}
	return nil
}
