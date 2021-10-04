package notification

import (
	"github.com/slack-go/slack"
	"os"
	"secure-pipeline-poc/app/config"
	"strings"
)

const channel = "secure-pipeline"

var APIURL = slack.APIURL

func Notify(messages []string) {
	if messages == nil {
		return
	}
	token := os.Getenv(config.SlackToken)
	client := slack.New(token, slack.OptionAPIURL(APIURL))

	text := strings.Join(messages, "\n\n")

	opt := slack.MsgOptionText(text, true)
	_, _, _, err := client.SendMessage(channel, opt)
	if err != nil {
		panic("Slack couldn't send a message!")
	}
}
