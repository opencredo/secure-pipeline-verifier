package notification

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"os"
	"secure-pipeline-poc/app/config"
	"strings"
)

const (
	Channel = "secure-pipeline"

	InfoMessage = "INFO"
	WarningMessage = "WARNING"
	ErrorMessage = "ERROR"

	InfoColor = "good"
	WarningColor = "warning"
	ErrorColor = "danger"
)

var APIURL = slack.APIURL


type MsgNotification struct {
	Control string `json:"control"`
	Level   string `json:"level"`
	Msg     string `json:"msg"`
}

func Notify(policyEvaluation interface{}) {

	token := os.Getenv(config.SlackToken)
	client := slack.New(token, slack.OptionAPIURL(APIURL))

	msgNotification, err := fillNotificationStruct(policyEvaluation)
	if err != nil {
		fmt.Println("Error collecting policy evaluation", err.Error())
		return
	}

	err = sendMessage(msgNotification, client)
	if err != nil {
		panic("Slack couldn't send a policyEvaluation!")
	}
}

func sendMessage(message MsgNotification, client *slack.Client) error {
	var err error
	if strings.Contains(message.Level, InfoMessage) {
		err = SendInfo(message, client)
	} else if strings.Contains(message.Level, WarningMessage) {
		err = SendWarning(message, client)
	} else if strings.Contains(message.Level, ErrorMessage) {
		err = SendError(message, client)
	}

	return err
}

func SendInfo(message MsgNotification, client *slack.Client) (err error) {
	return funcName(message, withAttachment(message, InfoColor), client)
}

func SendWarning(message MsgNotification, client *slack.Client) (err error) {
	return funcName(message, withAttachment(message, WarningColor), client)
}

func SendError(message MsgNotification, client *slack.Client) (err error) {
	return funcName(message, withAttachment(message, ErrorColor), client)
}

func funcName(text MsgNotification, attachment slack.Attachment, client *slack.Client) error {
	_, _, err := client.PostMessage(
		Channel,
		slack.MsgOptionText(text.Control, false),
		slack.MsgOptionAttachments(attachment),
	)

	return err
}

func withAttachment(message MsgNotification, color string) slack.Attachment {
	return slack.Attachment {
		Pretext: message.Level,
		Text: message.Msg,
		Color: color,
	}
}

func fillNotificationStruct(mapEval interface{}) (MsgNotification, error) {
	marshal, _ := json.Marshal(mapEval)

	var msgNotification MsgNotification
	err := json.Unmarshal(marshal, &msgNotification)
	if err != nil {
		return MsgNotification{}, err
	}
	fmt.Println("Notification: ", msgNotification)
	return msgNotification, nil
}
