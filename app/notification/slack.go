package notification

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"os"
	"secure-pipeline-poc/app/config"
)

const (
	InfoMessage    = "INFO"
	WarningMessage = "WARNING"
	ErrorMessage   = "ERROR"

	InfoColor    = "good"
	WarningColor = "warning"
	ErrorColor   = "danger"
)

var APIURL = slack.APIURL

type MsgNotification struct {
	Channel string
	Control string `json:"control"`
	Level   string `json:"level"`
	Msg     string `json:"msg"`
}

func Notify(policyEvaluation interface{}, slackConfig config.Slack) {
	if slackConfig.Enabled {
		token := os.Getenv(config.SlackToken)
		client := slack.New(token, slack.OptionAPIURL(APIURL))

		msgNotification, err := fillNotificationStruct(policyEvaluation, slackConfig.Channel)
		if err != nil {
			fmt.Println("Error collecting policy evaluation", err.Error())
			return
		}

		err = sendMessage(msgNotification, client)
		if err != nil {
            panic(fmt.Sprintf("Slack couldn't send a policyEvaluation!\n Error: %v", err))
		}
	}
}

func sendMessage(message MsgNotification, client *slack.Client) error {
	var err error
	switch message.Level {
	case InfoMessage:
		err = sendInfo(message, client)
	case WarningMessage:
		err = sendWarning(message, client)
	case ErrorMessage:
		err = sendError(message, client)
	}

	return err
}

func sendInfo(message MsgNotification, client *slack.Client) (err error) {
	return send(message, withAttachment(message, InfoColor), client)
}

func sendWarning(message MsgNotification, client *slack.Client) (err error) {
	return send(message, withAttachment(message, WarningColor), client)
}

func sendError(message MsgNotification, client *slack.Client) (err error) {
	return send(message, withAttachment(message, ErrorColor), client)
}

func send(message MsgNotification, attachment slack.Attachment, client *slack.Client) error {
	_, _, err := client.PostMessage(
		message.Channel,
		slack.MsgOptionText(message.Control, false),
		slack.MsgOptionAttachments(attachment),
	)

	return err
}

func withAttachment(message MsgNotification, color string) slack.Attachment {
	return slack.Attachment{
		Pretext: message.Level,
		Text:    message.Msg,
		Color:   color,
	}
}

func fillNotificationStruct(mapEval interface{}, channel string) (MsgNotification, error) {
	marshal, _ := json.Marshal(mapEval)

	var msgNotification MsgNotification
	err := json.Unmarshal(marshal, &msgNotification)
	if err != nil {
		return MsgNotification{}, err
	}
	msgNotification.Channel = channel
	return msgNotification, nil
}
